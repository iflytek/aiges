package instance

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/xfyun/aiges/buffer"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/aiges/protocol"
	"github.com/xfyun/aiges/storage"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// 任务状态查询/更新交互协议
type queryData struct {
	TaskId     string `json:"task_id"`
	TaskType   string `json:"task_type"`
	TaskStatue string `json:"task_status"`
	ReqParam   string `json:"req_param"`
}

type queryReq struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Sid     string      `json:"sid"`
	Data    []queryData `json:"data"`
}

type updateReq struct {
	TaskId     string `json:"task_id"`
	EngineNode string `json:"engine_node"`
	PrevStatus string `json:"prev_task_status"`
	TaskStatus string `json:"task_status"`
	RespResult []byte `json:"resp_result"`
	EngineCost int64  `json:"engine_cost_time"`
}

type updateResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Data    int    `json:"data"`
}

// 非实时任务状态
type taskStatus string

const (
	nrtCreate   taskStatus = "1"
	nrtActive   taskStatus = "2"
	nrtSuccess  taskStatus = "3"
	nrtExecFail taskStatus = "4"
)
const (
	nrtUpdateSucc int = 0
	nrtUpdateErr  int = 50403 // 数据库更新失败错误码
)

type nrtOp string

const (
	nrtOpSuc    nrtOp = "nrtSuccess" // 请求正常
	nrtOpIgnore nrtOp = "nrtIgnore " // 请求降级
	nrtOpFail   nrtOp = "nrtFail"    // 请求失败
)

// 非实时数据异步下行协议
type jsonDown struct {
	ReqParam map[string]string `json:"reqParam"`
	Data     []byte            `json:"data"`
	SpanMeta string            `json:"spanMeta"`
	Sid      string            `json:"sid"`
	Ret      int               `json:"ret"`
	ErrDesc  string            `json:"errDesc"`
}

// 非实时任务状态查询
func nrtQuery(task string) (op nrtOp, status taskStatus, err error) {
	query := conf.NrtDBUrl + "?" + "task_id=" + task
	resp, err := http.Get(query)
	if err != nil {
		op = nrtOpIgnore
		return
	} else if resp.StatusCode != 200 {
		err = errors.New("http get fail with " + resp.Status)
		op = nrtOpIgnore
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var msg queryReq
	err = json.Unmarshal(data, &msg)
	if err == nil {
		if msg.Code != 0 {
			op = nrtOpIgnore
			err = errors.New("http response fail with " + strconv.Itoa(msg.Code))
		} else if len(msg.Data) == 0 {
			op = nrtOpIgnore
			err = errors.New("http response data nil")
		} else {
			status = taskStatus(msg.Data[0].TaskStatue)
			op = nrtOpSuc
		}
	} else {
		op = nrtOpIgnore
	}
	return
}

// 非实时任务状态更新
func nrtUpdate(task string, cost int64, result []byte, status taskStatus) (op nrtOp, err error) {
	var cur taskStatus
	op, cur, err = nrtQuery(task) // 吴义平组提供的数据库服务接口设计,待讨论
	switch op {
	case nrtOpIgnore, nrtOpFail:
		return
	default:
	}

	req := updateReq{
		TaskId:     task,
		EngineNode: svcAddr + ":" + svcPort,
		TaskStatus: string(status),
		PrevStatus: string(cur),
		EngineCost: cost,
		RespResult: result}
	body, err1 := json.Marshal(req)
	if err1 != nil {
		return nrtOpIgnore, err1
	}
	putReq, err2 := http.NewRequest("PUT", conf.NrtDBUrl, bytes.NewReader(body))
	if err2 != nil {
		return nrtOpIgnore, err2
	}
	putResp, err3 := http.DefaultClient.Do(putReq)
	if err3 != nil {
		return nrtOpIgnore, err3
	}
	data, err := ioutil.ReadAll(putResp.Body)
	putResp.Body.Close()
	var resp updateResp
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nrtOpIgnore, err
	}
	if resp.Code != 0 {
		switch resp.Code {
		case nrtUpdateSucc:
			return nrtOpSuc, err
		//case nrtUpdateErr:
		//	return nrtOpFail, errors.New("nrtUpdate fail with " + strconv.Itoa(resp.Code))
		default:
			return nrtOpIgnore, errors.New("nrtUpdate fail with " + strconv.Itoa(resp.Code))
		}
	}
	return nrtOpSuc, err
}

// 非实时任务确认&更新
func nrtCheck(inst *ServiceInst) (errNum int, errInfo error) {
	if len(inst.headers[nrtTask]) > 0 {
		inst.nrtTime = time.Now().UnixNano()
		op, err := nrtUpdate(inst.headers[nrtTask], 0, nil, nrtActive)
		switch op {
		case nrtOpIgnore:
			inst.nrtDegrade = true
			inst.debug.WithTag("downgrade", err.Error())
			inst.tool.Log.Errorw("nrt task http downgrade", "task", inst.headers[nrtTask], "err", err.Error())
		case nrtOpFail:
			errNum, errInfo = frame.AigesErrorNrtUpdate, err
			inst.debug.WithErrorTag(err.Error()).WithRetTag(strconv.Itoa(errNum))
		default:
		}
	}
	return
}

// 非实时上行数据补齐
func nrtDataFill(bufData *[]buffer.DataMeta) (errNum int, errInfo error) {
	// check if data need to download from http/s3.
	for k, _ := range *bufData {
		ds, exist := (*bufData)[k].Desc.Attribute[dataSrc]
		if exist {
			switch string(ds) {
			case dataHttp:
				url, _ := (*bufData)[k].Desc.Attribute[dataHttpUrl]
				if len(url) == 0 {
					return frame.AigesErrorInvalidData, errors.New("input invalid http url")
				}

				// download from http, return err if download fail
				(*bufData)[k].Data, errNum, errInfo = storage.HttpDownload(string(url))
				if errInfo != nil {
					return frame.AigesErrorInvalidData, errInfo
				}
			case dataS3:
				access, _ := (*bufData)[k].Desc.Attribute[dataS3Access]
				secret, _ := (*bufData)[k].Desc.Attribute[dataS3Secret]
				endpoint, _ := (*bufData)[k].Desc.Attribute[dataS3Ep]
				bucket, _ := (*bufData)[k].Desc.Attribute[dataS3Bucket]
				key, _ := (*bufData)[k].Desc.Attribute[dataS3Key]
				if len(access) == 0 || len(secret) == 0 || len(endpoint) == 0 ||
					len(bucket) == 0 || len(key) == 0 {
					return frame.AigesErrorInvalidData, errors.New("input invalid s3 tags")
				}

				// TODO download from s3, return err if download fail
			case dataClient:
				// nothing to do
			default:
				return frame.AigesErrorInvalidData, errors.New("input invalid data source")
			}
		}
	}
	return
}

// 非实时下行数据推送
func dataPostBack(inst *ServiceInst, data *[]buffer.DataMeta, code int, err error) {
	var dataDown protocol.LoaderOutput
	var status protocol.LoaderOutput_RespStatus
	if data != nil {
		for _, rslt := range *data {
			var respRlt []byte
			if rslt.Data != nil {
				switch rslt.DataType {
				case buffer.DataText:
					respRlt = []byte(rslt.Data.(string))
				case buffer.DataAudio, buffer.DataImage, buffer.DataVideo:
					respRlt = rslt.Data.([]byte)
					// TODO 异步回传富媒体分离方案
				}
			}
			result := protocol.Payload{Meta: rslt.Desc, Data: respRlt}
			dataDown.Pl = append(dataDown.Pl, &result)
			status = protocol.LoaderOutput_RespStatus(rslt.Status)

		}
	}
	dataDown.Code = int32(code)
	if err != nil {
		dataDown.Err = err.Error()
	}
	dataDown.Status = status

	var jd jsonDown
	output, errMsl := proto.Marshal(&dataDown)
	jd.Sid = inst.instHdl
	jd.ReqParam = inst.params
	jd.SpanMeta = inst.spanMeta
	jd.Data = output
	jd.Ret = 0
	if errMsl != nil {
		inst.tool.Log.Errorw("dataPostBack proto marshal fail", "sid", inst.instHdl, "data", string(output), "err", errMsl.Error())
		jd.Ret = frame.AigesErrorPbMarshal
		jd.ErrDesc = frame.ErrorPbMarshal.Error()
	}
	jdata, merr := json.Marshal(jd)
	if merr != nil {
		inst.tool.Log.Errorw("dataPostBack json marshal fail", "sid", inst.instHdl, "jsondata", jd, "err", merr.Error())
	}

	// 更新数据库状态
	if !inst.nrtDegrade {
		nrtCost := (time.Now().UnixNano() - inst.nrtTime) / (1000 * 1000)
		nrtStatus := nrtSuccess
		if err != nil {
			nrtStatus = nrtExecFail
		}
		op, err := nrtUpdate(inst.headers[nrtTask], nrtCost, output, nrtStatus)
		switch op {
		case nrtOpIgnore:
			inst.nrtDegrade = true
			inst.debug.WithTag("degrade", err.Error())
			inst.tool.Log.Errorw("nrt task http downgrade", "task", inst.headers[nrtTask], "err", err.Error())
		case nrtOpFail:
			inst.debug.WithErrorTag(err.Error()).WithRetTag(strconv.Itoa(frame.AigesErrorNrtUpdate))
			inst.tool.Log.Errorw("nrt task http fail", "task", inst.headers[nrtTask], "err", err.Error())
			return
		default:
		}
	}

	// rabbitMQ 队列存储
	perr := storage.RabPublish(jdata, conf.RabbitRetry)
	if perr != nil {
		inst.tool.Log.Errorw("dataPostBack produce rmq message fail", "sid", inst.instHdl, "err", perr.Error())
	}
	inst.tool.Log.Debugw("dataPostBack produce rmq message", "sid", inst.instHdl, "body", jdata)
	return
}
