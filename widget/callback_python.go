package widget

import "C"
import (
	"errors"
	"github.com/xfyun/aiges/grpc/proto"
	"github.com/xfyun/aiges/grpc/shared"
	"github.com/xfyun/aiges/instance"
)

func engineCreateCallBackPy(in *proto.Response) (Code C.int) {
	// 异步回调接口,回调结果数据写缓冲区;
	var resp instance.ActMsg
	last := false

	if in.GetRet() != 0 {
		resp.AsyncErr = errors.New(enginePythonError(int(in.GetRet())))
		resp.AsyncCode = int(in.GetRet())
	} else {
		resp.DeliverData = make([]instance.DataMeta, 0, 1)
		for _, d := range in.List {
			var ele instance.DataMeta
			ele.DataId = d.GetKey()
			ele.DataStatus = int(d.GetStatus())
			ele.DataType = int(d.GetType())
			ele.Data = d.GetData()
			ele.DataDesc = d.GetDesc()
			resp.DeliverData = append(resp.DeliverData, ele)
			if ele.DataStatus == shared.DataEnd || ele.DataStatus == shared.DataOnce {
				// TODO 适配多数据流结果状态规则
				last = true
			}
		}
	}

	respChan, err := instance.QueryChan(in.GetTag())
	if err == nil {
		respChan <- resp
		if last || in.GetRet() != 0 {
			_ = instance.FreeChan(in.Tag)
		}
	} else {
		Code = C.int(-1)
	}
	return
}
