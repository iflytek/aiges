//go:build linux && cgo
// +build linux,cgo

package widget

/*
 #include "./widget/wrapper.h"
 #include <stdlib.h>
*/
import "C"
import (
	"errors"
	"github.com/xfyun/aiges/instance"
	"unsafe"
)

//export engineCreateCallBack
func engineCreateCallBack(hdl unsafe.Pointer, output C.pDataList, ret C.int) (Code C.int) {
	// 异步回调接口,回调结果数据写缓冲区;
	var resp instance.ActMsg
	last := false
	if ret != 0 {
		resp.AsyncErr = errors.New(engineError(int(ret)))
		resp.AsyncCode = int(ret)
	} else {
		resp.DeliverData = make([]instance.DataMeta, 0, 1)
		for output != nil {
			var ele instance.DataMeta
			ele.DataDesc = make(map[string]string)
			ele.DataId = C.GoString((*output).key)
			ele.DataType = int((*output)._type)
			ele.DataStatus = int((*output).status)
			pDesc := (*output).desc
			for pDesc != nil {
				ele.DataDesc[C.GoString((*pDesc).key)] = C.GoString((*pDesc).value)
				pDesc = (*pDesc).next
			}
			ele.Data = C.GoBytes(unsafe.Pointer((*output).data), C.int((*output).len))
			resp.DeliverData = append(resp.DeliverData, ele)
			output = (*output).next
			if ele.DataStatus == int(C.DataEnd) || ele.DataStatus == int(C.DataOnce) {
				// TODO 适配多数据流结果状态规则
				last = true
			}
		}
	}

	var usrTag string
	usrTag = C.GoString((*C.char)(hdl))
	cbchan, err := instance.QueryChan(usrTag)
	if err == nil {
		cbchan <- resp
		if last || ret != 0 {
			_ = instance.FreeChan(usrTag)
		}
	} else {
		Code = C.int(-1)
	}
	//	C.free(hdl)
	return
}
