/*
* @file	consistencyhashlb.go
* @brief  native 一致性哈希策略
*
* @author	jianjiang
* @version	1.0
* @date		2019.4
 */
package xsf

import (
	finder "github.com/xfyun/finder-go/common"
	"stathat.com/c/consistent"
	"sync"
)

type consistencyHashLB struct {
	sync.RWMutex

	sd          *serviceDiscovery
	hashCircles map[string]*consistent.Consistent //svc
}

func newConsistencyHashLB(o *conOption) *consistencyHashLB {
	conHashLb := new(consistencyHashLB)
	conHashLb.hashCircles = make(map[string]*consistent.Consistent)
	conHashLb.sd = newServiceDiscovery(o.fm)
	conHashLb.sd.registerCallBackFunc(conHashLb.updateHashCircle)
	return conHashLb
}

func (c *consistencyHashLB) updateHashCircle(svc string, notifyType string, instance []*finder.ServiceInstance) {
	if notifyType == string(finder.INSTANCEADDED) {
		for _, addrUnit := range instance {
			c.RLock()
			hashCircle, hashCircleOk := c.hashCircles[svc]
			c.RUnlock()
			if hashCircleOk {
				hashCircle.Add(addrUnit.Addr)
			} else {
				c.Lock()
				c.hashCircles[svc] = func() *consistent.Consistent { cycle := consistent.New(); cycle.Add(addrUnit.Addr); return cycle }()
				c.Unlock()
			}
		}
	} else if notifyType == string(finder.INSTANCEREMOVE) {
		for _, addrUnit := range instance {
			c.RLock()
			hashCircle, hashCircleOk := c.hashCircles[svc]
			c.RUnlock()
			if hashCircleOk {
				hashCircle.Remove(addrUnit.Addr)
			}
		}
	}
}

func (c *consistencyHashLB) Find(param *LBParams) ([]string, []string, error) {
	_, e := c.sd.findAll(param.version, param.svc, param.logId, param.log)
	if e != nil {
		return nil, nil, e
	}
	hashKey := param.hashKey
	c.RLock()
	hashCircle, hashCircleOk := c.hashCircles[param.svc]
	c.RUnlock()
	if !hashCircleOk {
		return nil, nil, EINAILIDSVC
	}
	addr, _ := hashCircle.Get(hashKey)
	if addr == "" {
		return []string{}, nil, nil
	}
	return []string{addr}, nil, nil
}
