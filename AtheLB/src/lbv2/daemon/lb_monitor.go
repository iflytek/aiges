package daemon

import (
	"encoding/json"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"sync"
	"time"
)

type monitorHandle interface {
	StrategyInst
	toMonitorWareHouse(*montorWareHouse)
	traversal() *monitorSvc
}

var monitorWareHouseInst montorWareHouse

const defaultMonitorInterval = time.Minute

type monitorSubSvcItem struct {
	Addr      string
	Timestamp int64 //纳秒,节点最近一次上报的时间
	BestInst  int64
	IdleInst  int64
	TotalInst int64
}
type monitorAddr struct {
	AddrMap map[string]*monitorSubSvcItem //K:addr
}

type monitorSubSvc struct {
	SubSvcMap map[string]*monitorAddr //K:sms
}

type monitorSvc struct {
	SvcMap map[string]*monitorSubSvc //K:iat
}

type montorWareHouse struct {
	AuthData     *monitorSvc
	AuthDataRwMu sync.RWMutex
	Svc          string //大业务名，如iat
	Defsub       string //兜底路由
	Ttl          int64  //纳秒,定时清除无效节点的时间间隔
	Threshold    int64
	RmqTopic     string
	RmqGroup     string

	MonitorInterval int64 //监控数据的更新时间,单位毫秒，缺省一分钟

	toolbox *xsf.ToolBox

	handle monitorHandle
}

func (m *montorWareHouse) update() {

	m.AuthDataRwMu.Lock()
	defer m.AuthDataRwMu.Unlock()
	m.AuthData = m.handle.traversal()

}
func (m *montorWareHouse) loop() {
	monitorInterval := defaultMonitorInterval
	if 0 != m.MonitorInterval {
		monitorInterval = time.Millisecond * time.Duration(m.MonitorInterval)
	}

	ticker := time.NewTicker(monitorInterval)

	for {
		select {
		case <-ticker.C:
			{
				m.toolbox.Log.Debugw("about to update monitor data", "monitorInterval", monitorInterval)
				m.update()
			}
		}
	}
}

/*
{
  "SubSvcMap": {
    "xox": {
      "AddrMap": null
    },
    "xux": {
      "AddrMap": {
        "x.x.x.x:oooo": {
          "Addr": "x.x.x.x:oooo",
          "Timestamp": 1542005397984167886,
          "BestInst": 0,
          "IdleInst": 0,
          "TotalInst": 0
        },
        "x.x.x.x:yyyy": {
          "Addr": "x.x.x.x:yyyy",
          "Timestamp": 1542005398175404814,
          "BestInst": 0,
          "IdleInst": 0,
          "TotalInst": 0
        }
      }
    }
  }
}
*/
func (m *montorWareHouse) query(svc, subsvc string, addr ...string) (rst string) {
	m.AuthDataRwMu.RLock()
	defer m.AuthDataRwMu.RUnlock()
	//什么都不传
	if "" != svc && "" != subsvc && 0 != len(addr) {
		rstByte, _ := json.Marshal(m)
		return string(rstByte)
	}

	//只传 svc
	if "" != svc && "" == subsvc && 0 == len(addr) {
		if nil == m.AuthData {
			return ""
		}
		svcDetail, svcDetailOk := m.AuthData.SvcMap[svc]

		m.toolbox.Log.Infow("fn:query", "svc", svc, "subsvc", subsvc, "addr", addr,
			"svcDetail", svcDetail, "svcDetailOk", svcDetailOk)

		if !svcDetailOk {
			return fmt.Sprintf("no svc:%v data", svc)
		}

		rstByte, _ := json.Marshal(svcDetail)
		return string(rstByte)
	}

	//只传 svc+subsvc
	if "" != svc && "" != subsvc && 0 == len(addr) {
		if nil == m.AuthData {
			return ""
		}
		svcDetail, svcDetailOk := m.AuthData.SvcMap[svc]
		if !svcDetailOk {
			return fmt.Sprintf("no svc:%v data", svc)
		}

		subSvcDetail, subSvcDetailOk := svcDetail.SubSvcMap[subsvc]
		if !subSvcDetailOk {
			return fmt.Sprintf("no svc:%v subsvc:%v data", svc, subsvc)
		}

		rstByte, _ := json.Marshal(subSvcDetail)
		return string(rstByte)

	}

	//只传 svc+svc+addr
	if "" != svc && "" != subsvc && 0 != len(addr) {
		if nil == m.AuthData {
			return ""
		}
		svcDetail, svcDetailOk := m.AuthData.SvcMap[svc]
		if !svcDetailOk {
			return fmt.Sprintf("no svc:%v data", svc)
		}

		subSvcDetail, subSvcDetailOk := svcDetail.SubSvcMap[subsvc]
		if !subSvcDetailOk {
			return fmt.Sprintf("no svc:%v subsvc:%v data", svc, subsvc)
		}

		rstOral := make(map[string]interface{})
		for _, v := range addr {
			subSvcAddrDetail, subSvcAddrDetailOk := subSvcDetail.AddrMap[v]
			if !subSvcAddrDetailOk {
				rstOral[v] = fmt.Sprintf("no svc:%v subsvc:%v addr:%v data", svc, subsvc, v)
				continue
			}
			rstByte, _ := json.Marshal(subSvcAddrDetail)
			rstOral[v] = string(rstByte)
		}
		rstByte, _ := json.Marshal(rstOral)

		return string(rstByte)
	}

	rstTmp, _ := json.Marshal(m)

	return string(rstTmp)
}
func (m *montorWareHouse) String() string {
	js, _ := json.Marshal(m)
	return string(js)
}
