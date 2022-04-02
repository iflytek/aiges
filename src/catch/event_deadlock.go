package catch

/**
此文件主要用来检测cgo进程卡死
设计文档见https://git.iflytek.com/kjsheng/notes/-/tree/master/aipaas%E6%8E%92%E9%9A%9C
*/

import "C"
import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type callType int

const (
	Begin          callType = 1
	End            callType = -1
	MaxReportTimes int      = 5
)

type callInfo struct {
	id    string //每一次调用的标识
	time  time.Time
	cType callType //1 标识请求刚进入 -1标识请求结束
	cycle int      //生命周期
}

var (
	currentReportTimes int
	jobTicker          *time.Ticker
	InfoChannel        chan callInfo
	controlFlag        chan bool
	infoMap            map[string]callInfo
	infoLock           sync.Mutex
	DeadLockSids       []string
)

type DeadlockEventHandle struct {
}

var deadLockHandler *DeadlockEventHandle

func (hd *DeadlockEventHandle) Occur() {
	if !switchOn {
		return
	}
	catchLog.Errorf("catch deadlock event")
	if MaxReportTimes < currentReportTimes {
		catchLog.Errorf("catch deadlock times bigger than setting. ignore")
		return
	}
	stack1, stack2 := hd.collectStack(event.pid)
	hd.reportEvent(stack1, stack2)
	currentReportTimes++
}
func (hd *DeadlockEventHandle) collectStack(pid string) (cStack string, goStack []byte) {
	cStack = GetPstack(pid)
	buf := make([]byte, 4<<20)
	n := runtime.Stack(buf, true)
	if n < (4 << 20) {
		return cStack, buf[:n]
	} else {
		catchLog.Errorw("catch deadlock. too large goroutine stack. only storage 4MB")
		return cStack, buf
	}
}
func (hd *DeadlockEventHandle) reportEvent(cStack string, goStack []byte) {
	dmg := dumpMsg{Cstack: cStack, Gostack: string(goStack), Type: "deadlock", Pid: event.pid, Host: event.host}
	dmg.Timestamp = time.Now().String()
	for _, sid := range DeadLockSids {
		dmg.Sessions = append(dmg.Sessions, sid)
	}

	msg, err := json.Marshal(dmg)
	if err != nil {
		catchLog.Errorw("dump json marshal fail with", "err", err.Error())
		return
	}

	dumpFile := dumpDir + "/deadlock_" + event.serviceId + "_" + event.pid + fmt.Sprintf("_%d", time.Now().Unix())
	file, err := os.Create(dumpFile)
	if err != nil {
		catchLog.Errorw("catch create dump file fail", "error", err.Error())
	} else {
		file.Write(msg)
		file.Write([]byte("\n")) // FileBeat采集依赖
		file.Close()
	}

	// TODO 当前仅支持crash事件
	notify("deadlock", dumpFile)
}

func GenerateUuid() string {
	uid := uuid.Must(uuid.NewV1(), nil)
	id := strings.Replace(uid.String(), "-", "", -1)
	return id
}

func CallCgo(uuid string, ctype callType) {
	if !switchOn {
		return
	}
	catchLog.Debugw("catch deadlock. call cgo", "ctype", ctype, "uuid", uuid)
	InfoChannel <- callInfo{
		id:    uuid,
		time:  time.Now(),
		cType: ctype,
	}
}

/*
一直处理 对任务信息预处理 写入map
*/
func DeadlockPretreatmentJob() {
	go func() {
		catchLog.Debugw("catch deadlock.  start deadlock pretreatment job")
		for {
			select {
			case val := <-InfoChannel:
				infoLock.Lock()
				if val.cType == Begin {
					infoMap[val.id] = val
				} else if val.cType == End {
					delete(infoMap, val.id)
				} else {
				}
				infoLock.Unlock()
			case <-controlFlag:
				catchLog.Debugw("catch deadlock.  stop deadlock pretreatment job")
				return
			}
		}
	}()
}

/*
定时任务 卡死检测
*/
func CornDeadlockDetectJob() {
	go func() {
		catchLog.Debugw("catch deadlock.  start deadlock detect job")
		for {
			select {
			case <-jobTicker.C:
				infoLock.Lock()
				for key, val := range infoMap {
					catchLog.Debugw("catch deadlock. ", "key", key, "value", val)
					if val.cycle >= 2 {
						if instMgrCallBack != nil {
							reqs := instMgrCallBack()
							for _, v := range reqs {
								DeadLockSids = append(DeadLockSids, v.Sid)
							}
						}
					} else {
						val.cycle += 1
						infoMap[key] = val
					}
				}
				if len(DeadLockSids) > 0 {
					catchLog.Debugw("catch deadlock.", "sids", DeadLockSids)
					deadLockHandler.Occur()
					DeadLockSids = nil
				}
				infoLock.Unlock()
			case <-controlFlag:
				catchLog.Debugw("catch deadlock.  stop deadlock detect job")
				return
			}
		}
	}()
}

func DeadLockDetectInit(interval int) {
	catchLog.Debugw("catch deadlock init")
	InfoChannel = make(chan callInfo, 100)
	infoMap = make(map[string]callInfo)
	controlFlag = make(chan bool)
	jobTicker = time.NewTicker(time.Second * time.Duration(interval))
	deadLockHandler = new(DeadlockEventHandle)
	CornDeadlockDetectJob()
	DeadlockPretreatmentJob()
}
func DeadLockDetectFini() {
	catchLog.Debugw("catch deadlock finish")
	close(InfoChannel)
	jobTicker.Stop()
	controlFlag <- true
	controlFlag <- true
	close(controlFlag)
}
