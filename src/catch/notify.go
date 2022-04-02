package catch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xfyun/aiges/conf"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 事件;
var event struct {
	pid          string
	host         string
	cloudId      string
	serviceId    string
	podName      string
	eventAPI     string
	podZone      string
	eventContext string
}

type dbgvar struct {
	name  string
	value *string
}

var dbgvars = []dbgvar{
	{"EVENT_CLOUD_ID", &event.cloudId},
	{"EVENT_SERVICE_ID", &event.serviceId},
	{"EVENT_POD_NAME", &event.podName},
	{"EVENT_PUSH_API", &event.eventAPI},
	{"EVENT_ZONE", &event.podZone},
	{"EVENT_CONTEXT", &event.eventContext},
}

func init() {
	for _, v := range dbgvars {
		*v.value = os.Getenv(v.name)
	}
	event.pid = strconv.Itoa(os.Getpid())
}

type ntfMeta struct {
	LogFile string `json:"file"`
	Pid     string `json:"pid"`
	Host    string `json:"host"`
	Time    string `json:"time"`
}

// 上报消息
type notification struct {
	Type      string  `json:"type"`
	CloudId   string  `json:"cloud_id"`
	Area      string  `json:"area"`
	ServiceId string  `json:"service_id"`
	Pod       string  `json:"pod"`
	Context   string  `json:"context"`
	MetaData  ntfMeta `json:"metadata"`
}

const retryLimit = 3

func notify(type_ string, logfile string) {
	ntf := notification{type_, event.cloudId, event.podZone, event.serviceId,
		event.podName, event.eventContext, ntfMeta{LogFile: logfile, Pid: event.pid, Host: conf.CatchSvcIP}}
	ntf.MetaData.Time = time.Now().String()
	msg, err := json.Marshal(ntf)
	if err != nil {
		if type_ == "start" {
			fmt.Println("loader event notify, json marshal fail with err", err.Error())
		} else {
			catchLog.Errorw("loader event notify, json marshal fail with", err, err.Error())
		}
		return
	}
	endFlag := make(chan bool)
	code := 0
	for i := 0; i < retryLimit; i++ {
		go func() {
			defer func() { endFlag <- true }()
			var resp *http.Response
			buf := bytes.NewBuffer(msg)
			resp, err = http.Post(event.eventAPI, "application/json", buf)
			if err != nil {
				code = -1
				if type_ == "start" {
					fmt.Println("loader event notify, http post with err", err.Error())
				} else {
					catchLog.Errorw("loader event notify, http post fail with", "err", err.Error())
				}
				return
			} else if resp.StatusCode != 200 {
				code = -1
				_ = resp.Body.Close()
				if type_ == "start" {
					fmt.Println("loader event notify, http response fail with err", resp.Status)
				} else {
					catchLog.Errorw("loader event notify, http response fail with", "err", resp.Status)
				}
				return
			} else {
				code = 0
				_ = resp.Body.Close()
				return
			}
		}()
		select {
		case <-time.After(time.Duration(5) * time.Second):
			fmt.Println("http post event fail time out > " + strconv.Itoa(conf.ResPerTimeout) + " ms")
			continue
		case <-endFlag:
			if code != 0 {
				continue
			}
			return
		}
	}
	if type_ == "start" {
		fmt.Println("loader post event end")
	} else {
		catchLog.Debugw("loader post event end")
	}
	// 应用层响应消息,忽略
}
