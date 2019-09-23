/*
	服务自检模块
	1. 提供服务预热
	2. 提供健康检查
	3. 提供优雅下线(done)
*/
package service

type sigClose struct {
	service *EngService
}

func (sc *sigClose) Closeout() {
	// xrpc 保证服务下线优雅性;
	sc.service.aiInst.FInit()
}
