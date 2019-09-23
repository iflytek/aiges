package server

import (
	"errors"
	"sync"
)

const MAX_DEPTH uint64 = 10
const DEFAULT_DEPTH uint64 = 5

const (
	//STATUS_ORIGINAL 代表原始
	STATUS_ORIGINAL uint32 = 0
	//STATUS_STARTING 代表正在启动
	STATUS_STARTING uint32 = 1
	//STATUS_STARTED 代表已启动
	STATUS_STARTED uint32 = 2
	//STATUS_STOPPING 代表正在停止
	STATUS_STOPPING uint32 = 3
	//STATUS_STOPPED 代表已停止
	STATUS_STOPPED uint32 = 4
)

//checkStatus 用于状态检测
//参数currentStatus 代表当前状态
//参数wantedStatus 代表想要的状态
func checkStatus(currentStatus uint32, wantedStatus uint32, lock sync.Locker) (err error) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	switch currentStatus {
	case STATUS_STARTING:
		err = errors.New("the pipeline is being opened")
	case STATUS_STOPPING:
		err = errors.New("the pipeline is being stopped")
	}
	if err != nil {
		return
	}
	switch wantedStatus {
	case STATUS_STARTING:
		if currentStatus == STATUS_STARTED {
			err = errors.New("the pipeline has been opened")
			return
		}
		if currentStatus == STATUS_STOPPED {
			err = errors.New("the pipeline has been closed")
			return
		}
	case STATUS_STOPPING:
		if currentStatus == STATUS_ORIGINAL {
			err = errors.New("the pipeline not yet opened")
			return
		}
		if currentStatus == STATUS_STOPPED {
			err = errors.New("the pipeline has been closed")
			return
		}
	}
	return
}
