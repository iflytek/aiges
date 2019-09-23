package instance

import (
	"bytes"
	"conf"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type queryData struct {
	TaskId string `json:"task_id"`
	TaskType string `json:"task_type"`
	TaskStatue string `json:"task_status"`
	ReqParam string `json:"req_param"`
}

type queryReq struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Sid string `json:"sid"`
	Data []queryData `json:"data"`
}

type updateReq struct {
	TaskId string `json:"task_id"`
	EngineNode string `json:"engine_node"`
	PrevStatus string `json:"prev_task_status"`
	TaskStatus string `json:"task_status"`
	RespResult []byte `json:"resp_result"`
	EngineCost int64 `json:"engine_cost_time"`
}

type updateResp struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Sid string `json:"sid"`
	Data int `json:"data"`
}

type taskStatus string
const (
	nrtCreate taskStatus = "1"
	nrtActive taskStatus = "2"
	nrtSuccess taskStatus = "3"
	nrtExecFail	taskStatus = "4"
)


// 非实时任务状态查询
func nrtQuery(task string) (status taskStatus, err error) {
	query := conf.NrtDBUrl + "?" + "task_id="+task
	resp, err := http.Get(query)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = errors.New("http get fail with " + resp.Status)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var msg queryReq
	err = json.Unmarshal(data, &msg)
	if err == nil && len(msg.Data) > 0 {
		status = taskStatus(msg.Data[0].TaskStatue)
	}
	return
}

// 非实时任务状态更新
func nrtUpdate(task string, cost int64, result []byte, status taskStatus) (err error) {
	prev,_ := nrtQuery(task)	// 数据库服务提供的接口设计,待讨论
	req := updateReq{
		TaskId:task,
		EngineNode:svcAddr + ":" + svcPort,
		TaskStatus:string(status),
		PrevStatus:string(prev),
		EngineCost:cost,	// 单位ms;
		RespResult:result}
	body, err := json.Marshal(req)
	if err == nil {
		putReq,_ := http.NewRequest("PUT", conf.NrtDBUrl, bytes.NewReader(body))
		putResp, err := http.DefaultClient.Do(putReq)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(putResp.Body)
		putResp.Body.Close()
		var resp updateResp
		err = json.Unmarshal(data, &resp)
		if err != nil{
			return  err
		}else if resp.Code != 0 {
			return errors.New("nrtUpdate fail with" + strconv.Itoa(resp.Code))
		}
	}
	return
}