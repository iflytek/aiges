/*
* @file	cleaner.go
* @brief  不活动地址清除
*         定期清除不活动的地址
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"sync"
	"time"
)

var (
	precision  = 5 * time.Second
	updateChan = 1024
)
// clData 记录服务地址以及访问时间的元结构
type clData struct {
	// addr 服务地址
	addr string

	// tm 访问时间,单位秒
	tm int
}

// clNode 维护服务地址访问记录的链表节点结构，双向循环链表
type clNode struct {
	// clData 节点数据
	td *clData

	// next 后节点
	next *clNode

	// 头节点
	pre *clNode
}

// cList 服务地址维护的双向链表，表头非空，非协程安全。
// 未选用包自带的链表，是由于实测过包的链表和个人实现的链表效率对比，从效率考虑采用自己造轮子
// todo: 补充测试数据
type clList struct {
	head *clNode
}

// init 初始化链表节点，如果cn不为空，则列为头节点
func (cl *clList) init(cn *clNode) {
	cl.head = cn
	if cl.head != nil {
		cl.head.next = nil
		cl.head.pre = nil
	}

}

// pushFront 在前头插入节点，如果链表为空，则先初始化
func (cl *clList) pushFront(cn *clNode) {
	if nil == cl.head {
		cl.init(cn)
	} else if nil != cl.head {
		cl.head.pre = cn
		cn.next = cl.head
		cn.pre = nil
		cl.head = cn
	}
}

// delNode 删除链表节点
func (cl *clList) delNode(cn *clNode) {
	//pre := cn.pre
	if cl.head == cn {
		cl.head = nil
		return
	}
	if nil != cn.next {
		cn.next.pre = cn.pre
	}
	cn.pre.next = cn.next
}

// moveToHead 把链表节点移动到表头
func (cl *clList) moveToHead(cn *clNode) {
	//pre := cn.pre
	if cl.head == cn {
		return
	}
	cn.pre.next = cn.next
	if nil != cn.next {
		cn.next.pre = cn.pre
	}
	cl.head.pre = cn
	cn.next = cl.head
	cn.pre = nil
	cl.head = cn
}

// addrVist 记录addr的访问时间
type addrVist struct {
	addr string
	vtm  int
}

// cleaner 清扫者对象
type cleaner struct {
	// 维护 l的线程同步
	sm sync.Mutex

	// l  访问者列表，便于移动和扫描
	l *clList

	// m 用于存放l节点的索引，便于查找
	m map[string]*clNode

	// t 定时器扫描l，超出过期时间的节点
	t *time.Timer

	// tm 定时器时间
	tm time.Duration

	// run  定时器启停标志
	run bool

	// lb 负载均衡句柄
	lb *loadBalance

	//当前时间,单位秒
	now int

	// addrVist 对象池
	avPool *sync.Pool

	log *utils.Logger

	//
	updateChan chan *addrVist
}

// newCleaner 创建cleaner对象，和sdk生命周期一致
func newCleaner(tm time.Duration, lb *loadBalance) *cleaner {
	c := new(cleaner)
	c.m = make(map[string]*clNode)
	c.run = true
	c.tm = tm
	c.log = lb.log
	//c.t = time.AfterFunc(tm,c.clean)
	c.t = time.AfterFunc(precision, c.updateTime)
	c.lb = lb
	c.l = &clList{}
	c.avPool = &sync.Pool{New: func() interface{} { return new(addrVist) }}
	c.updateChan = make(chan *addrVist, updateChan)
	go c.runner()
	return c
}

// update 更新addr的访问记录
func (c *cleaner) update(addr string) {

	vist := c.avPool.Get().(*addrVist)
	vist.addr = addr
	vist.vtm = c.now
	//todo 可能会阻塞
	c.updateChan <- vist

}

// runner 异步update访问时间到结构中
func (c *cleaner) runner() {

	tc := time.NewTicker(c.tm)
	for {
		logId := logSidGeneratorInst.GenerateSid("runner")
		select {
		case v := <-c.updateChan:

			c.log.Infow(
				"recv data from c.updateChan",
				"logId", logId, "addr", v.addr)
			//	c.sm.Lock()
			//	defer c.sm.Unlock()
			node, ok := c.m[v.addr]
			if ok {
				//  找到
				node.td.tm = c.now
				c.l.moveToHead(node)
			} else {
				n := clNode{td: &clData{addr: v.addr, tm: v.vtm}}
				c.m[v.addr] = &n
				c.l.pushFront(&n)
			}
			c.avPool.Put(v)
		case <-tc.C:
			c.log.Infow("recv data from tc.C")
			c.cleanWioutLock(logId)
		}
	}

}
func (c *cleaner) cleanWioutLock(logId string) {
	c.log.Infow("enter cleanWioutLock", "logId", logId)
	tmp := make([]string, 0, 10)
	if c.run {
		now := c.now

		c.log.Infow("cleanWioutLock", "logId", logId, "c.run", c.run)
		//	c.sm.Lock()
		tc := c.l.head
		for ; nil != tc; tc = tc.next {
			c.log.Infow("finding addrs who timeout", "logId", logId)
			// 发现过期addr
			//  if  now.Sub(tc.td.tm) > c.tm {
			if (now - tc.td.tm) > (int)(c.tm/time.Second) {
				c.log.Infow("found addrs who timeout", "logId", logId, "addr", tc.td.addr)
				c.l.delNode(tc)
				delete(c.m, tc.td.addr)
				tmp = append(tmp, tc.td.addr)
			} else {
				break
			}
		}
		//	c.sm.Unlock()

		for _, sv := range tmp {
			if len(sv) > 0 {
				c.lb.remove(sv)

			}
		}

		//	c.t = time.AfterFunc(c.tm,c.clean)
	}
}

// clean 定时任务，扫出过期地址
func (c *cleaner) clean() {
	/*
		这个函数好像没调用
	*/

	tmp := make([]string, 0, 10)
	if c.run {
		now := c.now

		c.sm.Lock()
		tc := c.l.head
		for ; nil != tc; tc = tc.next {

			// 发现过期addr
			//  if  now.Sub(tc.td.tm) > c.tm {
			if (now - tc.td.tm) > (int)(c.tm/time.Second) {
				c.l.delNode(tc)
				delete(c.m, tc.td.addr)
				tmp = append(tmp, tc.td.addr)
			} else {
				break
			}
		}
		c.sm.Unlock()

		for _, sv := range tmp {
			if len(sv) > 0 {
				c.lb.remove(sv)

			}
		}

		c.t = time.AfterFunc(c.tm, c.clean)
	}

}

// updateTime 更新时间,时间精度决定了过期地址清空的时间精度
func (c *cleaner) updateTime() {
	c.now = int(time.Now().UnixNano() / int64(time.Second))
	time.AfterFunc(precision, c.updateTime)
}

// stop停止过期维护
func (c *cleaner) stop() {
	c.run = false
	c.t.Stop()
}
