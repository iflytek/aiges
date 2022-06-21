package analy

import (
	"bufio"
	"encoding/json"
	"fmt"
	_var "github.com/xfyun/aiges/xtest/var"
	"github.com/xfyun/xsf/utils"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"
)

type direction int
type DataStatus int
type SessStatus int

const (
	UP   direction = 1
	DOWN direction = 2

	DataBegin    DataStatus = 0
	DataContinue DataStatus = 1
	DataEnd      DataStatus = 2
	DataTotal    DataStatus = 3

	SessBegin    SessStatus = 0
	SessContinue SessStatus = 1
	SessEnd      SessStatus = 2
	SessOnce     SessStatus = 3

	outputPerfFile   = "perf.txt"
	outputRecordFile = "perfReqRecord.csv"
	outputPerfImg    = "perf.jpg"
)

/*
xtest 性能检测工具
*/
type callDetail struct {
	ID       string     //uuid
	Handle   string     //会话模式时的hdl
	Tm       time.Time  //时间戳
	dataStat DataStatus //数据状态 ，0,1,2,3
	sessStat SessStatus //会话状态,0,1,2,3
	Dire     direction  //输入 还是输出
	ErrCode  int
	ErrInfo  string
}

type performance struct {
	Max         float32 `json:"max"`
	Min         float32 `json:"min"`
	FailRate    float32 `json:"failRate"`
	SuccessRate float32 `json:"successRate"`
	//平均值95 99线
	Delay95      float32 `json:"delay95"`
	Delay99      float32 `json:"delay99"`
	DelayAverage float32 `json:"delayAverage"`
	//首结果95 99线
	DelayFirstMin     float32 `json:"delayFirstMin"`
	DelayFirstMax     float32 `json:"delayFirstMax"`
	DelayFirst95      float32 `json:"delayFirst95"`
	DelayFirst99      float32 `json:"delayFirst99"`
	DelayFirstAverage float32 `json:"delayFirstAverage"`
	//尾结果95 99线
	DelayLastMin     float32 `json:"delayLatMin"`
	DelayLastMax     float32 `json:"delayLatMax"`
	DelayLast95      float32 `json:"delayLast95"`
	DelayLast99      float32 `json:"delayLast99"`
	DelayLastAverage float32 `json:"delayLastAverage"`
}

type singlePerfCost struct {
	id        string
	cost      float32
	firstCost float32 //首个结果耗时
	lastCost  float32 //最后一个结果耗时
}

type errMsg struct {
	ErrInfo string `json:"errInfo"`
	Handle  string `json:"handle"`
}

type PerfModule struct {
	idx            int
	collectChan    chan callDetail
	mtx            sync.Mutex
	control        chan bool
	correctReqPath map[string][]callDetail //正确的请求路径图

	errReqRecord map[int][]errMsg //错误的请求记录

	correctReqCost []singlePerfCost //正确的请求花费的时间记录

	perf performance //性能结果

	reqRecordFile *os.File

	Log *utils.Logger
}

var Perf *PerfModule

func (pf *PerfModule) Start() (err error) {
	pf.Log.Debugw("perf start")
	pf.control = make(chan bool)
	pf.collectChan = make(chan callDetail, 10000)
	pf.errReqRecord = make(map[int][]errMsg)
	pf.correctReqPath = make(map[string][]callDetail)
	pf.reqRecordFile, err = os.OpenFile(outputRecordFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	go func() {
		pf.collect()
	}()
	return nil
}
func (pf *PerfModule) Stop() {
	_ = pf.reqRecordFile.Close()
	pf.calcDelay()
	pf.dump()
	pf.control <- true
	close(pf.collectChan)
	pf.Log.Debugw("perf calc quit")
}

func (pf *PerfModule) Record(id, handle string, stat DataStatus, stat2 SessStatus, dire direction, errCode int, errInfo string) {
	if pf == nil {
		return
	}
	tmp := time.Now()
	pf.collectChan <- callDetail{
		id, handle, tmp, stat, stat2, dire, errCode, errInfo,
	}
}

func (pf *PerfModule) collect() {
	pf.Log.Debugw("perf start collect")
	for {
		select {
		case detail := <-pf.collectChan:
			pf.Log.Debugw("perf collect", "detail ", detail)
			pf.mtx.Lock()
			if !_var.ReqMode {
				if detail.ErrInfo != "" {
					pf.errReqRecord[detail.ErrCode] = append(pf.errReqRecord[detail.ErrCode], errMsg{Handle: detail.ID, ErrInfo: detail.ErrInfo})
				} else {
					pf.correctReqPath[detail.ID] = append(pf.correctReqPath[detail.ID], detail)
					if detail.Dire == DOWN {
						pf.pretreatment(detail.ID)
						delete(pf.correctReqPath, detail.ID)
					}
				}
			} else {
				//创建会话req首帧
				if detail.sessStat == SessBegin && detail.Dire == UP {
					pf.correctReqPath[detail.ID] = make([]callDetail, 1)
					pf.correctReqPath[detail.ID][0] = detail
					pf.Log.Debugw("get first segment up", "id", detail.ID, "tm", time.Now(), "detail", detail)
				} else if detail.sessStat == SessBegin && detail.Dire == DOWN { //创建会话resp首帧
					pf.Log.Debugw("get first segment down", "id", detail.ID, "tm", time.Now(), "detail", detail)
					tmp := pf.correctReqPath[detail.ID][0] //获取创建会话的第一帧 因为没有handle，所以通过id来拿
					tmp.Handle = detail.Handle
					delete(pf.correctReqPath, detail.ID)
					if detail.ErrInfo != "" {
						pf.errReqRecord[detail.ErrCode] = append(pf.errReqRecord[detail.ErrCode], errMsg{Handle: detail.Handle, ErrInfo: detail.ErrInfo})
					} else {
						pf.correctReqPath[detail.Handle] = append(pf.correctReqPath[detail.Handle], tmp)
						pf.correctReqPath[detail.Handle] = append(pf.correctReqPath[detail.Handle], detail)
					}
				} else {
					if detail.ErrInfo != "" {
						pf.errReqRecord[detail.ErrCode] = append(pf.errReqRecord[detail.ErrCode], errMsg{Handle: detail.Handle, ErrInfo: detail.ErrInfo})
						delete(pf.correctReqPath, detail.Handle)
					} else {
						pf.correctReqPath[detail.Handle] = append(pf.correctReqPath[detail.Handle], detail)
						if detail.Dire == DOWN && detail.dataStat == DataEnd {
							pf.pretreatment(detail.Handle)
							delete(pf.correctReqPath, detail.Handle)
						}
					}
				}
			}
			pf.mtx.Unlock()
		case <-pf.control:
			return
		}
	}
}

func (pf *PerfModule) pretreatment(id string) {
	if !_var.ReqMode {
		var begin, end time.Time
		for _, record := range pf.correctReqPath[id] {
			if record.Dire == UP {
				begin = record.Tm
			} else {
				end = record.Tm
			}
		}
		pf.idx += 1
		_, _ = pf.reqRecordFile.WriteString(fmt.Sprintf("id:%s,cost:%f,begin:%s,end:%s\n",
			id, float32(end.Sub(begin).Microseconds())/1000, begin, end))
		if pf.idx%500 == 0 {
			_ = pf.reqRecordFile.Sync()
		}
	} else {
		var begin, end time.Time
		var firstReq, firstRlt time.Time
		var lastReq, lastRlt time.Time
		knowUp := false
		knowDown := false
		for _, record := range pf.correctReqPath[id] {
			if record.Dire == UP {
				if record.dataStat == DataBegin && !knowUp {
					begin = record.Tm
					firstReq = record.Tm
					knowUp = true
				} else if record.dataStat == DataEnd {
					lastReq = record.Tm
				} else {

				}
			} else {
				if record.dataStat == DataBegin && !knowDown {
					firstRlt = record.Tm
					knowDown = true
				} else if record.dataStat == DataEnd {
					lastRlt = record.Tm
					end = record.Tm
				} else {

				}
			}
		}
		pf.idx += 1
		if !knowUp {
			begin = lastReq
			firstReq = lastReq
		}
		if !knowDown {
			firstRlt = end
		}
		_, _ = pf.reqRecordFile.WriteString(fmt.Sprintf("id:%s,cost:%f,firstCost:%f,lastCost:%f,begin:%s,end:%s,"+
			"firstReq:%s,firstRlt:%s,lastReq:%s,lastRlt:%s\n",
			id, float32(end.Sub(begin).Microseconds())/1000,
			float32(firstRlt.Sub(firstReq).Microseconds())/1000,
			float32(lastRlt.Sub(lastReq).Microseconds())/1000,
			begin, end, firstReq, firstRlt, lastReq, lastRlt))
		if pf.idx%500 == 0 {
			_ = pf.reqRecordFile.Sync()
		}
	}
}

func (pf *PerfModule) loadRecord() error {
	pf.Log.Debugw("perf load record")
	loadFile, err := os.Open(outputRecordFile)
	if err != nil {
		pf.Log.Errorf("perf failed to load record file. %s", err.Error())
		return err
	}
	defer loadFile.Close()
	br := bufio.NewReader(loadFile)
	for {
		record, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if !_var.ReqMode {
			reg := regexp.MustCompile("id:(.*?),cost:(.*?),begin:(.*?),end:(.*?)")
			parts := reg.FindStringSubmatch(string(record))
			if len(parts) != 5 {
				pf.Log.Errorw("perf failed to regex record", "record", string(record))
				return nil
			}
			cost, _ := strconv.ParseFloat(parts[2], 32)
			pf.correctReqCost = append(pf.correctReqCost, singlePerfCost{
				id:   parts[1],
				cost: float32(cost),
			})
		} else {
			reg := regexp.MustCompile("id:(.*?),cost:(.*?),firstCost:(.*?),lastCost:(.*?),begin:(.*?),end:(.*?)," +
				"firstReq:(.*?),firstRlt:(.*?),lastReq:(.*?),lastRlt:(.*?)")
			parts := reg.FindStringSubmatch(string(record))
			if len(parts) != 11 {
				pf.Log.Errorw("perf failed to regex record", "record", string(record))
				return nil
			}
			cost, _ := strconv.Atoi(parts[2])
			firstCost, _ := strconv.ParseFloat(parts[3], 32)
			lastCost, _ := strconv.ParseFloat(parts[4], 32)
			pf.correctReqCost = append(pf.correctReqCost, singlePerfCost{
				id:        parts[1],
				cost:      float32(cost),
				firstCost: float32(firstCost),
				lastCost:  float32(lastCost),
			})
		}
	}
	return nil
}

func (pf *PerfModule) calcDelay() {
	pf.Log.Debugw("perf calc delay start")
	err := pf.loadRecord()
	if err != nil {
		return
	}
	if len(pf.correctReqCost) == 0 {
		pf.Log.Debugw("perf correct number of request is zero")
		return
	}
	//非会话模式或者会话模式仅保留第一帧到最后一帧的性能结果
	if !_var.ReqMode || _var.PerfLevel == 0 {
		for _, val := range pf.correctReqCost {
			pf.Log.Debugf("perf info,id:%s,cost:%f\n", val.id, val.cost)
		}
		pf.perf.Min, pf.perf.Max, pf.perf.DelayAverage, pf.perf.Delay95, pf.perf.Delay99 =
			pf.anallyArray(func(costs []singlePerfCost) []float32 {
				var tmp []float32
				for _, v := range costs {
					tmp = append(tmp, v.cost)
				}
				return tmp
			}(pf.correctReqCost))
	} else if _var.ReqMode && _var.PerfLevel == 1 {
		pf.perf.DelayFirstMin, pf.perf.DelayFirstMax, pf.perf.DelayFirstAverage, pf.perf.DelayFirst95, pf.perf.DelayFirst99 =
			pf.anallyArray(func(costs []singlePerfCost) []float32 {
				var tmp []float32
				for _, v := range costs {
					tmp = append(tmp, v.firstCost)
				}
				return tmp
			}(pf.correctReqCost))
	} else if _var.ReqMode && _var.PerfLevel == 2 {
		pf.perf.DelayLastMin, pf.perf.DelayLastMax, pf.perf.DelayLastAverage, pf.perf.DelayLast95, pf.perf.DelayLast99 =
			pf.anallyArray(func(costs []singlePerfCost) []float32 {
				var tmp []float32
				for _, v := range costs {
					tmp = append(tmp, v.lastCost)
				}
				return tmp
			}(pf.correctReqCost))
	} else {
	}

	var errCount int
	for k, _ := range pf.errReqRecord {
		errCount += len(pf.errReqRecord[k])
	}
	pf.Log.Debugw("perf calc", "correctNum", len(pf.correctReqCost), "errCount", errCount)
	pf.perf.SuccessRate = float32(len(pf.correctReqCost)) / float32(errCount+len(pf.correctReqCost))
	pf.perf.FailRate = float32(math.Abs(float64(1 - pf.perf.SuccessRate)))
}

func (pf *PerfModule) dump() {
	file, err := os.OpenFile(outputPerfFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		pf.Log.Errorf("perf failed to dump performance data.%s ", err)
		return
	}
	defer file.Close()
	val, _ := json.Marshal(pf.perf)
	_, _ = file.WriteString("perf result:\n ")
	_, _ = file.Write(val)
	_, _ = file.WriteString("\n")

	val, _ = json.Marshal(pf.errReqRecord)
	_, _ = file.WriteString("perf err record: \n")
	_, _ = file.Write(val)
	_, _ = file.WriteString("\n")
}

func (pf *PerfModule) anallyArray(data []float32) (min, max, average, aver95, aver99 float32) {
	if len(data) == 1 {
		return data[0], data[0], data[0], data[0], data[0]
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})
	min = data[0]
	max = data[len(data)-1]
	average, aver95, aver99 = func(costs []float32) (a, b, c float32) {
		var totalCost float32
		for _, val := range costs {
			totalCost += val
		}
		pf.Log.Debugw("anallyArray ", "costs ", costs)
		return totalCost / float32(len(costs)),
			costs[(len(costs)-1)*95/100],
			costs[(len(costs)-1)*99/100]
	}(data)
	return
}
