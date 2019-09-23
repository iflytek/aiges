package daemon

type GetInPutOpt func(*GetInPut)

func withGetSvc(svc string) GetInPutOpt {
	return func(in *GetInPut) {
		in.svc = svc
	}
}
func withGetSubSvc(subSvc string) GetInPutOpt {
	return func(in *GetInPut) {
		in.subSvc = subSvc
	}
}
func withGetNBest(nBest int64) GetInPutOpt {
	return func(in *GetInPut) {
		in.nBest = nBest
	}
}
func withGetAll(all bool) GetInPutOpt {
	return func(in *GetInPut) {
		in.all = all
	}
}
func withGetUid(uid int64) GetInPutOpt {
	return func(in *GetInPut) {
		in.uid = uid
	}
}

func withGetExParam(exParam string) GetInPutOpt {
	return func(in *GetInPut) {
		in.exParam = exParam
	}
}
func withGetSid(sid string) GetInPutOpt {
	return func(in *GetInPut) {
		in.sid = sid
	}
}

func withGetPersonalized(personalized bool) GetInPutOpt {
	return func(in *GetInPut) {
		in.personalized = personalized
	}
}

type GetInPut struct {
	sid          string //lb内部数据聚合所用
	all          bool   //返回所有的节点数据，测试用
	uid          int64
	svc          string
	subSvc       string
	exParam      string
	nBest        int64
	personalized bool //标识此次查询是否为个性化相关
}
