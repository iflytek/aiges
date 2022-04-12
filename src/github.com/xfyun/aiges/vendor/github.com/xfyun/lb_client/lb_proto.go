package lb_client

import (
	"github.com/golang/protobuf/proto"
	lbClientPb "github.com/xfyun/lbClientPb"
)

//封装注册的proto信息
func marshalLbLoginMsg(svc string, totalInst, idleInst, bestInst int32, param map[string]string) (data []byte, err error) {
	if param == nil {
		param = make(map[string]string)
	}

	lbPb := &lbClientPb.LbReport{
		Addr:      proto.String(svc),
		TotalInst: proto.Int32(totalInst),
		IdleInst:  proto.Int32(idleInst),
		BestInst:  proto.Int32(bestInst),
		Param:     param,
	}
	data, err = proto.Marshal(lbPb)
	if err != nil {
		return
	}
	return
}

//封装更新的proto信息
func marshalLbUpdateMsg(data []byte, totalInst, idleInst, bestInst int32) (updateData []byte, err error) {
	lbPb := &lbClientPb.LbReport{}
	if err = proto.Unmarshal(data, lbPb); err != nil {
		return
	}

	*lbPb.TotalInst = totalInst
	*lbPb.IdleInst = idleInst
	*lbPb.BestInst = bestInst

	updateData, err = proto.Marshal(lbPb)
	if err != nil {
		return
	}
	return
}
