package catch

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

/**
捕获 启动异常时的相关堆栈和信息
*/
type StartFailedHandle struct {
}

func (hd *StartFailedHandle) Occur() {
	if !switchOn {
		return
	}
	fmt.Println("catch start failed event")
	stack1, stack2 := hd.collectStack(event.pid)
	hd.reportEvent(stack1, stack2)
}

func (hd *StartFailedHandle) collectStack(pid string) (cStack string, goStack []byte) {
	cStack = GetPstack(pid)
	buf := make([]byte, 4<<20)
	n := runtime.Stack(buf, true)
	if n < (4 << 20) {
		return cStack, buf[:n]
	} else {
		fmt.Println("catch startfailed. too large goroutine stack. only storage 4MB")
		return cStack, buf
	}
}

func (hd *StartFailedHandle) reportEvent(cStack string, goStack []byte) {
	dmg := dumpMsg{Cstack: cStack, Gostack: string(goStack), Type: "start", Pid: event.pid, Host: event.host}
	dmg.Timestamp = time.Now().String()
	msg, err := json.Marshal(dmg)
	if err != nil {
		fmt.Println("dump json marshal fail with,err", err.Error())
		return
	}

	dumpFile := dumpDir + "/start_" + event.serviceId + "_" + event.pid + fmt.Sprintf("_%d", time.Now().Unix())
	file, err := os.Create(dumpFile)
	if err != nil {
		fmt.Println("catch create dump file fail error", err.Error())
	} else {
		file.Write(msg)
		file.Write([]byte("\n")) // FileBeat采集依赖
		file.Close()
	}

	// TODO 当前仅支持crash事件
	notify("start", dumpFile)
}
