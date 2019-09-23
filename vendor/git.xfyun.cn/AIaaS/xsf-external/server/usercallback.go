package xsf

import (
	"fmt"
	"strings"
	"sync"
)

var callBackList map[string]UserCallBack
var callBackListRwMu sync.RWMutex

const (
	UserHighPriority   = "__0"
	UserNormalPriority = "__1"
	UserLowPriority    = "__2"
)

func init() {
	callBackListRwMu.Lock()
	callBackList = make(map[string]UserCallBack)
	callBackListRwMu.Unlock()
}

type UserCallBack interface {
	Exec()
}

func AddUserCallBack(handle string, task UserCallBack) error {
	return addUserCallBack(handle, task, UserNormalPriority)
}
func AddUserCallBackWithPriority(handle string, task UserCallBack, priority string) error {
	return addUserCallBack(handle, task, priority)
}

func addUserCallBack(handle string, task UserCallBack, priority string) error {
	callBackListRwMu.Lock()
	defer callBackListRwMu.Unlock()
	if _, ok := callBackList[handle+priority]; ok {
		return fmt.Errorf("handle:%v already stored", handle)
	}
	callBackList[handle+priority] = task

	return nil
}

func dealUserCallBack() {
	/*
		按批次执行，简单的遍历，有很大的优化空间，后续优化
	*/
	callBackListRwMu.RLock()

	/*
		执行高优先级
	*/
	loggerStd.Printf("deal with UserHighPriority\n")
	for handle, task := range callBackList {
		if strings.Contains(handle, UserHighPriority) {
			task.Exec()
		}
	}

	/*
		执行普通优先级
	*/
	loggerStd.Printf("deal with UserNormalPriority\n")
	for handle, task := range callBackList {
		if strings.Contains(handle, UserNormalPriority) {
			task.Exec()
		}
	}

	/*
		执行低优先级
	*/
	loggerStd.Printf("deal with UserLowPriority\n")
	for handle, task := range callBackList {
		if strings.Contains(handle, UserLowPriority) {
			task.Exec()
		}
	}

	callBackListRwMu.RUnlock()
}
