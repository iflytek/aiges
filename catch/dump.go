package catch

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type dumpMsg struct {
	Cstack    string   `json:"cstack"`
	Gostack   string   `json:"gostack"`
	Sessions  []string `json:"sid"`
	Type      string   `json:"type"`
	Pid       string   `json:"pid"`
	Timestamp string   `json:"timestamp"`
	Host      string   `json:"host"`
}

func dump(cStack []byte, goStack []byte, info error) {
	dmg := dumpMsg{Cstack: string(cStack), Gostack: string(goStack), Type: "crash", Pid: event.pid, Host: event.host}
	dmg.Timestamp = time.Now().String()
	if instMgrCallBack != nil {
		reqs := instMgrCallBack()
		for _, v := range reqs {
			dmg.Sessions = append(dmg.Sessions, v.Sid)
		}
	}

	msg, err := json.Marshal(dmg)
	if err != nil {
		fmt.Println("dump json marshal fail with ", err.Error())
		return
	}

	dumpFile := dumpDir + "/crash_" + event.serviceId + "_" + event.pid + fmt.Sprintf("_%d", time.Now().Unix())
	file, err := os.Create(dumpFile)
	if err != nil {
		catchLog.Errorw("catch create dump file fail", "error", err.Error())
	} else {
		file.Write(msg)
		file.Write([]byte("\n")) // FileBeat采集依赖
		file.Close()
	}

	// TODO 当前仅支持crash事件
	notify("crash", dumpFile)
}
