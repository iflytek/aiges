/*
* @file	discover.go
* @brief  服务发现
*         根据服务名管理的一批可用的服务的addr
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

// todo: 需要优化本地配置中读取目标地址的逻辑
// todo: proxy mode、 weight、loadbalance mode 需要生效

package xsf

import (
	"errors"
	"git.xfyun.cn/AIaaS/finder-go/common"
	"strings"
	"sync"
	"sync/atomic"

	//"fmt"

	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

var (
	TESTADDR []string
	//addrs_map map[string]*cList
	addrsMap sync.Map
)

// 用于测试
func init() {
	TESTADDR = make([]string, 0, 10)
	TESTADDR = append(TESTADDR, "127.0.0.1")
	//addrs_map = make(map[string]*cList)
}

// setTeatAddr 设置配置文件中获取的服务的地址
func setTeatAddr(addr string) error {
	TESTADDR = strings.Split(addr, ",")
	if len(TESTADDR) <= 0 {
		return errors.New("Invalid addr cfg")
	}
	for _, a := range TESTADDR {
		t := strings.Split(a, "@")
		if len(t) <= 1 {
			return errors.New("Invalid addr cfg")
		}
		sub := t[0]
		var cl cList
		t1 := strings.Split(t[1], ";")
		for _, a1 := range t1 {
			if len(a1) < 1 {
				return errors.New("Invalid addr cfg")
			}
			cl.insert(&node{addr: a1})

		}
		//addrs_map[sub] = &cl
		addrsMap.Store(sub, &cl)
	}
	return nil
}

//todo 获取全部地址的接口 getAll
// node 服务链表单节点的结构
type node struct {
	// addr 服务地址
	addr string

	next *node
	pre  *node
}

// cList 带游标的双向循环链表，非协程安全
type cList struct {
	cIx  int64 //当前的滚动游标
	cLen int64 //链表长度
	cur  *node
	head *node
}

// init 初始化链表
func (cl *cList) init(n *node) {
	cl.head = n
	cl.cur = cl.head
	if nil != n {
		cl.head.next = cl.head
		cl.head.pre = cl.head
	}
}

// insert 插入节点
func (cl *cList) insert(n *node) {
	if nil == cl.head && nil != n {

		cl.init(n)
	} else if nil != n { //todo 貌似有bug
		t := cl.head.next
		cl.head.next = n
		cl.head.next.next = t

		t.pre = n
		n.pre = cl.head
	}
	atomic.AddInt64(&cl.cLen, 1)
}

// next 获取下个节点
func (cl *cList) next() (*node) {
	if nil == cl.cur {
		return nil
	}

	cl.cur = cl.cur.next
	return cl.cur
}

// next 获取下个节点
func (cl *cList) nextInList(s int) []string {
	// todo::
	if nil == cl.head {
		return nil
	}

	if 0 == s {
		//此时获取全量地址，供hash策略使用
		addrs := make([]string, 0, s)
		tmp := cl.head
		for {
			addrs = append(addrs, tmp.addr)
			tmp = tmp.next

			if tmp == cl.head {
				return addrs
			}
		}
	}

	j := atomic.AddInt64(&cl.cIx, 1) % cl.cLen

	addrs := make([]string, 0, s)
	tmp := cl.head
	for i := int64(0); i < j; i++ {
		tmp = tmp.next
	}

	for sIx := 0; sIx < s; sIx++ {
		addrs = append(addrs, tmp.addr)
		tmp = tmp.next
	}
	return addrs
}

// delete 删除节点
func (cl *cList) delete(n *node) {
	if n == cl.head { //表头，
		if cl.head == cl.head.next { //且只有表头
			cl.head = nil
			cl.cur = nil

		} else {

			cl.head = cl.head.next
			cl.cur = cl.head
			pre := n.pre
			next := n.next
			pre.next = n.next
			next.pre = n.pre

		}
		atomic.AddInt64(&cl.cLen, -1)
		n = nil
		return
	}
	atomic.AddInt64(&cl.cLen, -1)
	pre := n.pre
	next := n.next
	pre.next = n.next
	next.pre = n.pre
	n = nil
}

// cMap 服务名和可用地址列表的映射关系
type cMap struct {
	addrs map[string]*node
	cl    *cList
	mu    sync.RWMutex
}

// insert 插入地址
func (m *cMap) insert(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n, ok := m.addrs[addr]
	if ok {
		return
	}
	n = &node{addr: addr}
	m.cl.insert(n)
	m.addrs[addr] = n

}

// setlist 设置list
func (m *cMap) setlist(cl *cList) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// todo:
	/*n ,ok := m.addrs[addr]
	if ok {
		n.w = w
		return
	}
	n = &node{addr:addr,w:w}*/
	m.cl = cl
	//m.addrs[addr] = n

}

// next 获取下一个地址
func (m *cMap) Next() (*node) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if nil == m.cl {
		return nil
	}
	return m.cl.next()
}

// next 获取下一个地址
func (m *cMap) NextInList(s int) ([]string, []string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if nil == m.cl {
		return nil, nil
	}
	return m.cl.nextInList(s), m.cl.nextInList(0)
}

// delete 删除一个地址
func (m *cMap) delete(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n, ok := m.addrs[addr]
	if ok {
		delete(m.addrs, addr)

		//todo:删除list，更新Node
		/*	if m.cl != nil {
				for i :=0 ;i < 11;i ++{
					fmt.Println( m.cl.cur.addr)
					m.cl.cur = m.cl.cur.next
				}

			}
			fmt.Println("delete")*/
		m.cl.delete(n)
		/*
			if m.cl != nil {
				fmt.Println("delete",m.cl.cur)

				for i :=0 ;i < 11;i ++{
					fmt.Println( m.cl.cur.addr)
					m.cl.cur = m.cl.cur.next
				}

			}*/
	}

}

type ServiceDiscoryCallBack = func(string, string, []*finder.ServiceInstance)

// service 表示一个服务对象地址池
type service struct {
	addrs *cMap
	//	ProxyMode       string
	//	LoadBalanceMode string
}

// newService 申请一个服务对象池
func newService() *service {

	ss := new(service)
	ss.addrs = new(cMap)
	ss.addrs.addrs = make(map[string]*node)
	ss.addrs.cl = new(cList)
	return ss
}

// serviceDiscovery 表示走配置中心存放的服务对象池管理器
type serviceDiscovery struct {
	finder *utils.FindManger

	// mss 服务名到地址列表的索引表，key: 服务名，value:服务对象地址池
	// 相当于map[string]service
	mss          *sync.Map //[string]*service
	callBackFunc []ServiceDiscoryCallBack
}

// newServiceDiscovery 申请一个服务对象管理器的对象
func newServiceDiscovery(finder *utils.FindManger) *serviceDiscovery {
	sd := new(serviceDiscovery)
	//sd.ss = make(map[string]*service)
	sd.mss = new(sync.Map)
	sd.finder = finder
	return sd
}

func (sd *serviceDiscovery) registerCallBackFunc(cbf ServiceDiscoryCallBack) {
	sd.callBackFunc = append(sd.callBackFunc, cbf)
}

// 服务发现注册的回调接口，服务实例和配置均发生了变更
func (sd *serviceDiscovery) OnServiceInstanceConfigChanged(
	name string,
	version string,
	instance string,
	config *finder.ServiceInstanceConfig) bool {
	loggerStd.Println("OnServiceInstanceConfigChanged")
	//	s,ok := sd.ss[name]
	s, ok := sd.mss.Load(name)
	if !ok {
		return false
	}
	if config.IsValid {
		loggerStd.Println("OnServiceInstanceConfigChanged insert adddr: ", instance)

		s.(*service).addrs.insert(instance)
	} else {
		loggerStd.Println("OnServiceInstanceConfigChanged delete adddr: ", instance)
		s.(*service).addrs.delete(instance)
	}

	return true
}

// 服务发现注册的回调接口，服务发现服务配置发生了变化
func (sd *serviceDiscovery) OnServiceConfigChanged(
	name string,
	version string,
	config *finder.ServiceConfig) bool {
	loggerStd.Println("OnServiceConfigChanged")

	return true
}

// 服务发现注册的回调接口，服务实例发生变更
func (sd *serviceDiscovery) OnServiceInstanceChanged(name string,
	version string,
	instances []*finder.ServiceInstanceChangedEvent) bool {
	loggerStd.Println("OnServiceInstanceChanged")

	s := sd.findService(name)
	if nil == s {
		return false
	}
	for _, v := range instances {
		if v.EventType == finder.INSTANCEADDED {
			for _, inst := range v.ServerList {
				//	fmt.Println("OnServiceInstanceChanged INSTANCEADDED",inst.Addr)
				loggerStd.Println("OnServiceInstanceChanged insert addrs:", inst.Addr)

				s.addrs.insert(inst.Addr)
			}
		} else if v.EventType == finder.INSTANCEREMOVE {
			for _, inst := range v.ServerList {
				//	fmt.Println("OnServiceInstanceChanged INSTANCEREMOVE",inst.Addr)
				loggerStd.Println("OnServiceInstanceChanged delete addrs:", inst.Addr)

				s.addrs.delete(inst.Addr)
			}
		}
		for _, cbf := range sd.callBackFunc {
			cbf(name, string(v.EventType), v.ServerList)
		}
	}
	return true
}
func (sd *serviceDiscovery) findAllService() map[string]*service {
	rst := make(map[string]*service)
	sd.mss.Range(func(key, value interface{}) bool {
		rst[key.(string)] = value.(*service)
		return true
	})
	return rst
}

// findService 根据服务名从缓存中查询服务地址池
func (sd *serviceDiscovery) findService(name string) *service {
	s, ok := sd.mss.Load(name)
	if ok {
		//fmt.Println("findService", s.(*service))

		return s.(*service)
	}
	//fmt.Println("findService not nil")

	return nil
}

// insertService 插入服务名->地址池的映射关系
func (sd *serviceDiscovery) insertService(name string, s *service) {
	//sd.ss[name] = s
	sd.mss.Store(name, s)
}

// findAll 根据服务名查询服务地址池，如果地址不存在，则在配置中心中查找可用服务。查询完成后，缓存到本地列表
func (sd *serviceDiscovery) findAll(version, name, logId string, log *Logger) (*service, error) {
	// todo: 函数太长，注意优化
	//fmt.Println("enter")
	ss := sd.findService(name)
	if nil != ss {
		log.Infow("found service from local cache", "name", name, "logId", logId)
		return ss, nil
	}
	log.Infow("can't found service from local cache", "name", name, "logId", logId)

	if nil == sd.finder {
		return sd.initFromAddrsMap(log, logId, name, ss)
	}

	ss = newService()
	//fmt.Println("in")
	log.Infow(
		"calling sd.finder.UseSrvAndSub",
		"logId", logId, "version", version, "name", name)

	srv, e := sd.finder.UseSrvAndSub(version, name, sd)

	if nil != e {
		log.Infow(
			"failed to call sd.finder.UseSrvAndSub",
			"logId", logId, "err", e)

		return nil, e
	}
	log.Infow(
		"success to call sd.finder.UseSrvAndSub",
		"logId", logId, "srv", srv, "err", e)

	if nil == srv {
		return nil, INVALIDSRV
	}

	//todo: 异常处理
	hasVal := false
	for _, l := range srv[name+"_"+version].ProviderList {
		if l.Config.IsValid {
			hasVal = true

			ss.addrs.insert(l.Addr)
		}
	}
	for _, cbf := range sd.callBackFunc {
		cbf(name, string(finder.INSTANCEADDED), srv[name+"_"+version].ProviderList)
	}
	log.Infow(
		"dealing hasVal",
		"logId", logId, "hasVal", hasVal)
	if !hasVal {
		return nil, INVALIDSRV
	}

	sd.insertService(name, ss)

	return ss, nil

}

func (sd *serviceDiscovery) initFromAddrsMap(log *Logger, logId string, name string, ss *service) (*service, error) {
	log.Infow("sd.finder is nil", "logId", logId, "name", name)
	//l, ok := addrs_map[name]
	l, ok := addrsMap.Load(name)
	if !ok {
		log.Infow(
			"failed take name from addrs_map.Load(name)",
			"logId", logId, "name", name)
		return nil, INVALIDFINDER
	}

	log.Infow(
		"success take name from addrs_map.Load(name),calling newService()",
		"logId", logId, "name", name)
	//	cursor := l.head
	ss = newService()
	//ss.addrs.insert(cursor.addr,cursor.w)

	ss.addrs.setlist(l.(*cList))

	//同步地址到topKLb
	var instances []*finder.ServiceInstance
	for _, addr := range l.(*cList).nextInList(0) {
		instances = append(instances, &finder.ServiceInstance{Addr: addr})
	}
	for _, cbf := range sd.callBackFunc {
		cbf(name, string(finder.INSTANCEADDED), instances)
	}

	sd.insertService(name, ss)

	return ss, nil

}
