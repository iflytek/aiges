/*
 *提供引擎上下线管理，节点信息更新以及个性化资源失效的更新
 */
package lb_client

import (
	"errors"
	"path"
	"time"

	"github.com/cihub/seelog"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	Err_Svc_Is_Null = errors.New("svc is null")
	engine_node     string //引擎节点(ip:port)
)

const (
	TRY_MAX_CNT       = 3 //最大尝试次数
	HEALTH_CHECK_TIME = 5 //zk健康检查时间,单位:秒
)

type LbClienter interface {
	Init(o ...LbCfgOpt) (err error)                                                             //初始化
	Login(svc string, totalInst, idleInst, bestInst int32, param map[string]string) (err error) //引擎节点上线
	LoginOut() (err error)                                                                      //引擎节点主动下线
	Upadate(totalInst, idleInst, bestInst int32) (err error)                                    //更新引擎节点信息
}

type LbClient struct {
	LbUtil
}

//初始化
func (lc *LbClient) Init(o ...LbCfgOpt) (err error) {
	err = lc.LbUtil.Init(o...)
	if err != nil {
		return err
	}
	return err
}

/*
 *引擎节点注册
 *svc-引擎节点(ip:port),totleInst-总实例数，idleInst-空闲实例数，bestInst-最优实例数，param-传递的参数
 */
func (lc *LbClient) Login(svc string, totalInst, idleInst, bestInst int32, param map[string]string) (err error) {
	if len(svc) == 0 {
		err = Err_Svc_Is_Null
	}
	engine_node = svc

	data, err := marshalLbLoginMsg(svc, totalInst, idleInst, bestInst, param)
	if err != nil {
		return
	}

	for _, subRouterType := range lc.LbOpt.SubRouterTypes {
		var suc bool = false
		for i := 0; i < TRY_MAX_CNT; i++ {
			err = lc.createAliveNode(svc, subRouterType, data)
			if err != nil {
				continue
			}
			suc = true
			break
		}
		if !suc {
			return
		}
	}

	//开启zk状态监控
	go lc.watchZkStatusAndProcess(svc, data)
	return
}

//引擎主动下线
func (lc *LbClient) LoginOut() (err error) {
	for _, subRouterType := range lc.LbOpt.SubRouterTypes {
		var suc bool = false
		for i := 0; i < TRY_MAX_CNT; i++ {
			err = lc.deleteAliveNode(engine_node, subRouterType)
			if err != nil {
				continue
			}
			suc = true
			break
		}
		if !suc {
			return
		}
	}

	return
}

/*
 *更新引擎节点信息
 *totleInst-总实例数，idleInst-空闲实例数，bestInst-最优实例数
 */
func (lc *LbClient) Upadate(totalInst, idleInst, bestInst int32) (err error) {
	for _, subRouterType := range lc.SubRouterTypes {
		svcAliveAddr := path.Join(lc.RouterTypeAbs, subRouterType, engine_node)
		svcData, _, err := lc.ZkConn.Get(svcAliveAddr)
		if err != nil {
			return err
		}

		updateSvcData, err := marshalLbUpdateMsg(svcData, totalInst, idleInst, bestInst)
		if err != nil {
			return err
		}

		//更新存活服务器目录下的该服务节点
		if _, err = lc.ZkConn.Set(svcAliveAddr, updateSvcData, -1); err != nil {
			return err
		}
	}
	return
}

//zk状态定时监控
func (lc *LbClient) watchZkStatusAndProcess(svc string, data []byte) {
	t1 := time.NewTimer(time.Second * time.Duration(HEALTH_CHECK_TIME))
	for {
		select {
		case <-t1.C:
			lc.zkStatusChangedProcess(svc, data)
			t1.Reset(time.Second * time.Duration(HEALTH_CHECK_TIME))
		}
	}
}

//zk状态改变做相应的处理
func (lc *LbClient) zkStatusChangedProcess(svc string, data []byte) {
	defer func() {
		if errErr := recover(); errErr != nil {
			seelog.Error("occur panic,err is:", errErr)
		}
	}()
	statusOk := lc.LbUtil.zkHealthCheck()
	if statusOk {
		for _, subRouterType := range lc.SubRouterTypes {
			nodeAddr := path.Join(lc.RouterTypeAbs, subRouterType, svc)
			exists, _, err := lc.ZkConn.Exists(nodeAddr)
			if err != nil {
				continue
			}
			if !exists {
				lc.ZkConn.Create(nodeAddr, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			}
		}
	} else {
		seelog.Critical("zk lost connect!!!")
	}
}
