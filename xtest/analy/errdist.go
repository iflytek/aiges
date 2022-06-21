package analy

import (
	_var "github.com/xfyun/aiges/xtest/var"
	"github.com/xfyun/xsf/utils"
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

type errInfo struct {
	errCode int
	errStr  error
}
type errDistAnalyser struct {
	errCnt   map[int]int64 // map[error]count 错误计数
	errDsc   map[int]error // 错误描述
	errTmp   []errInfo     // 临时存储区,用于channel满阻塞的极端场景;
	errMutex sync.Mutex
	errChan  chan errInfo // error
	swg      sync.WaitGroup
	log      *utils.Logger
}

func (eda *errDistAnalyser) Start(clen int, logger *utils.Logger) {
	eda.log = logger
	eda.errCnt = make(map[int]int64)
	eda.errDsc = make(map[int]error)
	eda.errTmp = make([]errInfo, 0, 10)
	eda.errChan = make(chan errInfo, clen)
	eda.swg.Add(1)
	go eda.count()
}

func (eda *errDistAnalyser) PushErr(code int, err error) {
	select {
	case eda.errChan <- errInfo{code, err}:
	default:
		// channel满阻塞,降级加锁写入临时存储区
		eda.errMutex.Lock()
		defer eda.errMutex.Unlock()
		eda.errTmp = append(eda.errTmp, errInfo{code, err})
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

		cnt, _ := eda.errCnt[err.errCode]
		eda.errCnt[err.errCode] = cnt + 1
		eda.errDsc[err.errCode] = err.errStr
	}

	// 临时存储区数据同步
	eda.errMutex.Lock()
	for _, v := range eda.errTmp {
		cnt, _ := eda.errCnt[v.errCode]
		eda.errCnt[v.errCode] = cnt + 1
		eda.errDsc[v.errCode] = v.errStr
	}
	eda.errMutex.Unlock()

	// dump Log to disk file;
	eda.dumpLog()
	eda.swg.Done()
}

func (eda *errDistAnalyser) dumpLog() {
	// 错误分布数据落盘;
	fi, err := os.OpenFile(_var.ErrAnaDst, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		eda.log.Errorw("error dist Log dump fail with open file", "err", err.Error(), "file", _var.ErrAnaDst)
		return
	}

	for eCode, eCount := range eda.errCnt {
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
