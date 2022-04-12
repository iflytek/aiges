package xsf

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var signalHandleStarted int64 //>0 stand for already start
var KillerCheckList map[string]Killer
var KillerMu sync.RWMutex
var KillerCheckListInit int64

const (
	killerFirstPriority  = "__0" //服务主动下线，停止上报，等待授权释放完毕或超时
	killerHighPriority   = "__1" //执行用户FINI逻辑
	killerNormalPriority = "__2" //停止grpc的监听行为，处理剩余请求
	killerLowPriority    = "__3" //日志落盘
	killerLastPriority   = "__4" //等待程序结束
)

func init() {
	if atomic.CompareAndSwapInt64(&KillerCheckListInit, 0, 1) {
		KillerCheckList = make(map[string]Killer)
	}
}

type Killer interface {
	Closeout()
}

func AddKillerCheck(handle string, task Killer) {
	addKillerCheck(killerNormalPriority, handle, task)
}

/*
	将优先级信息已双下划线为分割符，拼接于handle后面
*/
func addKillerCheck(priority string, handle string, task Killer) {
	KillerMu.Lock()
	if atomic.CompareAndSwapInt64(&KillerCheckListInit, 0, 1) {
		KillerCheckList = make(map[string]Killer)
	}
	KillerCheckList[handle+priority] = task
	KillerMu.Unlock()
}
func GracefulStop() {
	deathSentence()
}

//仅保grpc优雅退出和block部分运行重复执行，其它只能执行一次
var deathSentenceOnce1 sync.Once
var deathSentenceOnce2 sync.Once

func deathSentence() {

	/*
		按批次执行，简单的遍历，有很大的优化空间，后续优化
	*/
	KillerMu.RLock()

	deathSentenceOnce1.Do(func() {
		{
			/*
				最高优先级
			*/
			loggerStd.Printf("deal with killerFirstPriority\n")
			for handle, task := range KillerCheckList {
				if strings.Contains(handle, killerFirstPriority) {
					task.Closeout()
				}
			}

			/*
				执行高优先级
			*/
			loggerStd.Printf("deal with killerHighPriority\n")
			for handle, task := range KillerCheckList {
				if strings.Contains(handle, killerHighPriority) {
					task.Closeout()
				}
			}
		}
	})
	/*
		执行普通优先级
	*/
	loggerStd.Printf("deal with killerNormalPriority\n")
	for handle, task := range KillerCheckList {
		if strings.Contains(handle, killerNormalPriority) {
			task.Closeout()
		}
	}
	deathSentenceOnce2.Do(func() {
		{
			/*
				执行低优先级
			*/
			loggerStd.Printf("deal with killerLowPriority\n")
			for handle, task := range KillerCheckList {
				if strings.Contains(handle, killerLowPriority) {
					task.Closeout()
				}
			}
		}
	})

	/*
		最低优先级
	*/
	loggerStd.Printf("deal with killerLastPriority\n")
	for handle, task := range KillerCheckList {
		if strings.Contains(handle, killerLastPriority) {
			task.Closeout()
		}
	}

	KillerMu.RUnlock()
}
func signalHandle() {
	if !atomic.CompareAndSwapInt64(&signalHandleStarted, 0, 1) {
		return
	}

end:
	for {

		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGPIPE)
	loop:
		s := <-c
		switch s {
		case syscall.SIGTERM, syscall.SIGINT:
			{
				loggerStd.Printf("ts:%v,receive syscall.SIGTERM\n", time.Now())

				KillerMu.RLock()
				if len(KillerCheckList) == 1 && func() bool {
					_, ok := KillerCheckList["finder"+killerFirstPriority]
					return ok
				}() {
					loggerStd.Printf("KillerCheckList only contain finder,break.\n")
					signal.Stop(c)
					break end
				}
				KillerMu.RUnlock()

				deathSentence()
				//os.Exit(returnCode)
			}
		case syscall.SIGPIPE:
			{
				goto loop
			}
		}
	}
}
func SignalHandle() {
	signalHandle()
}
