package daemon

import "fmt"

type SetInPutOpt func(*SetInPut)

func withSetLive(live int) SetInPutOpt {
	return func(in *SetInPut) {
		in.live = live
	}
}
func withSetAddr(addr string) SetInPutOpt {
	return func(in *SetInPut) {
		in.addr = addr
	}
}
func withSetSvc(svc string) SetInPutOpt {
	return func(in *SetInPut) {
		in.svc = svc
	}
}
func withSetSubSvc(subSvc string) SetInPutOpt {
	return func(in *SetInPut) {
		in.subSvc = subSvc
	}
}
func withSetTotal(total int64) SetInPutOpt {
	return func(in *SetInPut) {
		in.total = total
	}
}
func withSetIdle(idle int64) SetInPutOpt {
	return func(in *SetInPut) {
		in.idle = idle
	}
}
func withSetBest(best int64) SetInPutOpt {
	return func(in *SetInPut) {
		in.best = best
	}
}
func withSetSid(sid string) SetInPutOpt {
	return func(in *SetInPut) {
		in.sid = sid
	}
}

type SetInPut struct {
	sid               string //lb内部数据聚合所用
	live              int
	addr              string
	svc               string
	subSvc            string
	total, idle, best int64
}

func (S *SetInPut) String() string {
	return fmt.Sprintf("live:%d,addr:%s,svc:%s,subsvc:%s,total:%d,idle:%d,best:%d",
		S.live, S.addr, S.svc, S.subSvc, S.total, S.idle, S.best)
}
