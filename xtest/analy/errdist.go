package analy

import (
	"git.iflytek.com/AIaaS/xsf/utils"
	"os"
	"strconv"
	"sync"
)

/*
用于成功率统计, 统计维度如下：
1. 会话/非会话请求成功率；
2. 错误码分布统计；
数据结构：map[int]int (map[error]count)
*/

var ErrAnalyser errDistAnalyser

type ErrInfo struct {
	ErrCode int
	ErrStr  error
}
type errDistAnalyser struct {
	errCnt   map[int]int64 // map[error]count 错误计数
	errDsc   map[int]error // 错误描述
	errTmp   []ErrInfo     // 临时存储区,用于channel满阻塞的极端场景;
	errMutex sync.Mutex
	errChan  chan ErrInfo // error
	swg      sync.WaitGroup
	log      *utils.Logger

	ErrAnaDst string
}

func (eda *errDistAnalyser) Start(clen int, logger *utils.Logger, errAnaDst string) {
	eda.log = logger
	eda.errCnt = make(map[int]int64)
	eda.errDsc = make(map[int]error)
	eda.errTmp = make([]ErrInfo, 0, 10)
	eda.errChan = make(chan ErrInfo, clen)
	eda.ErrAnaDst = errAnaDst
	eda.swg.Add(1)
	go eda.count()
}

func (eda *errDistAnalyser) PushErr(info ErrInfo) {
	select {
	case eda.errChan <- info:
	default:
		// channel满阻塞,降级加锁写入临时存储区
		eda.errMutex.Lock()
		defer eda.errMutex.Unlock()
		eda.errTmp = append(eda.errTmp, info)
	}
}

func (eda *errDistAnalyser) Stop() {
	close(eda.errChan)
	eda.swg.Wait()
}

func (eda *errDistAnalyser) count() {
	for {
		err, calive := <-eda.errChan
		if !calive {
			break
		}

		cnt, _ := eda.errCnt[err.ErrCode]
		eda.errCnt[err.ErrCode] = cnt + 1
		eda.errDsc[err.ErrCode] = err.ErrStr
	}

	// 临时存储区数据同步
	eda.errMutex.Lock()
	for _, v := range eda.errTmp {
		cnt, _ := eda.errCnt[v.ErrCode]
		eda.errCnt[v.ErrCode] = cnt + 1
		eda.errDsc[v.ErrCode] = v.ErrStr
	}
	eda.errMutex.Unlock()

	// dump Log to disk file;
	eda.dumpLog()
	eda.swg.Done()
}

func (eda *errDistAnalyser) dumpLog() {
	// 错误分布数据落盘;
	fi, err := os.OpenFile(eda.ErrAnaDst, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		eda.log.Errorw("error dist Log dump fail with open file", "err", err.Error(), "file", eda.ErrAnaDst)
		return
	}
	//delete(eda.errCnt, 0)
	for eCode, eCount := range eda.errCnt {
		//if eCode == 0 { // 正确请求，直接返回
		//	continue
		//}
		var edesc string
		if eda.errDsc[eCode] != nil {
			edesc = eda.errDsc[eCode].Error()
		}
		tmp := []byte(strconv.Itoa(eCode) + "(\"" + edesc + "\")" + ": ")
		tmp = append(tmp, []byte(strconv.Itoa(int(eCount)))...)
		tmp = append(tmp, byte('\n'))
		wlen, err := fi.Write(tmp)
		if err != nil {
			eda.log.Errorw("error dist Log dump fail with write ana", "err", err.Error(), "wlen", wlen)
			_ = fi.Close()
			return
		}
	}
	_ = fi.Close()
}
