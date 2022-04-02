/*
	catch包旨在加载器进程出现异常(panic/crash)时, 对触发异常的现场进行捕获及转移

当前支持：
	1. 异常触发go panic	:唯一定位异常请求
	2. 异常触发go signal	:暂无法唯一定位
	3. 异常触发c signal	:暂无法唯一定位
后续方向：
	1. 对异常signal捕获现场进行调研分析, 唯一定位异常请求
	2. 区分异常实例可重置场景及不可重置场景, 区分hook处理策略
*/
package catch

/*
#include <signal.h>
#include <stdio.h>
#include <execinfo.h>
#include <stdlib.h>
#include <string.h>

const int maxArr = 1024;
extern void goSigHandler(int sig, char** sym, int symLen);
// TODO c层堆栈可参考：runtime.printCgoTraceback && runtime.SetCgoTraceback
void cSigHandler(int sig){
	void* DumpArray[maxArr];
	char** symbols = NULL;
	char* failinfo = NULL;
	int nSize = backtrace(DumpArray, maxArr);
	if (nSize > 0) {
		symbols = backtrace_symbols(DumpArray, nSize);
	}else {
		failinfo = (char*) malloc(50);
		strcpy(failinfo, "catch sighandler backtrace fail\n");
		symbols = &failinfo;
	}

	goSigHandler(sig, symbols, nSize);
	if (failinfo != NULL){
		free(failinfo);
	}else if (symbols != NULL){
		free(symbols);
	}
	exit(0);
}

void sigRegister() {
	// register signal handler
	signal(SIGSEGV, cSigHandler);
	signal(SIGABRT, cSigHandler);
}
*/
import "C"
import (
	"errors"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var instMgrCallBack func() (reqDoubt []TagRequest)

var sigErrTable = map[syscall.Signal]error{
	syscall.SIGSEGV: errors.New("SIGSEGV: segmentation violation"),
	syscall.SIGABRT: errors.New("SIGABRT: abort"),
}

func signalHandle() {
	//	if switchOn {
	// 捕获c层signal
	C.sigRegister()
	// 捕获go层signal异常
	signal.Notify(sigChan, syscall.SIGSEGV)
	signal.Notify(sigChan, syscall.SIGABRT)

	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGSEGV, syscall.SIGABRT:
			dump(nil, debug.Stack(), sigErrTable[sig.(syscall.Signal)])
		default:
			catchLog.Errorw("catch SignalHook | receive unknown signal. ", "sig", sig)
		}
	}()
	//	}
}
