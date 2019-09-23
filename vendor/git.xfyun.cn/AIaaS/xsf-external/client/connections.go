/*
* @file	connPoolections.go
* @brief  提供rpc连接管理操作方法
*         提供连接创建、销毁、自动恢复等功能
*
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	"sync"
	"sync/atomic"
	"time"
)

//  connPoolNode 连接维护链表节点的结构
type connPoolNode struct {
	cc   *grpc.ClientConn
	next *connPoolNode
	pre  *connPoolNode
}

// connPoolMeta 接维护链表结构，双向循环链表。协程安全
type connPoolMeta struct {
	// 连接池当前大小
	size int
	// 连接池最大大小
	max int
	// 连接读缓冲区大小，默认GRCMRS
	rbuf int
	// 连接写缓冲区大小，默认GRCMWS
	wbuf       int
	maxReceive int

	//keepalive检查时间间隔
	keepaliveTime time.Duration

	//keepalive检查超时时间
	keepaliveTimeout time.Duration

	cur  *connPoolNode
	head *connPoolNode
	m    sync.RWMutex

	poolCnt int64
}

// init 链表初始化
func (cl *connPoolMeta) init(n *connPoolNode) {
	cl.head = n
	cl.cur = cl.head
	if n != nil {
		cl.head.next = cl.head
		cl.head.pre = cl.head
	}
}

// insert 链表插入
func (cl *connPoolMeta) insert(n *connPoolNode) bool {
	cl.m.Lock()
	defer cl.m.Unlock()
	if cl.size >= cl.max {
		return false
	}
	if nil == cl.head && nil != n {
		cl.init(n)
		cl.size += 1
		return true
	} else if nil != n {
		t := cl.head.next
		cl.head.next = n
		n.pre = cl.head
		n.next = t
		t.pre = n
		cl.size += 1
		return true
	}
	return false
}

// del 删除一个连接节点
// todo:性能待优化
func (cl *connPoolMeta) del(n *connPoolNode) bool {
	cl.m.Lock()
	defer cl.m.Unlock()
	if nil != cl.head {
		t := cl.head
		if cl.head == n {
			//只有1个节点
			if cl.head.next == cl.head {
				cl.head = nil
				cl.cur = nil
				cl.size -= 1
				n.cc.Close()

				return true
			}
			th := cl.head.pre
			cl.head = cl.head.next
			th.next = cl.head
			cl.head.pre = th
			cl.size -= 1
			n.cc.Close()

			return true
		}

		for t = t.next; t != cl.head; t = t.next {
			// 找到并删除
			if n == t {
				// 如果删除的是当前游标节点，重置游标
				if cl.cur == n {
					cl.cur = t.next
				}
				th := t.pre
				th.next = t.next
				t.next.pre = th
				cl.size -= 1
				n.cc.Close()

				return true
			}
		}
	}
	return false
}

// next 获取下个节点
func (cl *connPoolMeta) next() (*connPoolNode) {
	cl.m.RLock()
	defer cl.m.RUnlock()
	if 0 == cl.size {
		return nil
	}
	cl.cur = cl.cur.next
	return cl.cur
}

// alloc 分配一个可用连接，如果连接池不满，则创建新的连接
func (cm *connPoolMeta) alloc(addr string, ms time.Duration) (*grpc.ClientConn, error) {
	// 如果连接数小于最大数，重新分配连接
	if cm.size < cm.max {
		c, e := cm.newConnPool(addr, ms)
		if nil != c || nil != e {
			return c, e
		}
	}
	//  从池子中取一个链接
	return cm.get(addr, ms)
}
func (cm *connPoolMeta) fillConnPool(done chan error, addr string, ms time.Duration) {
	if atomic.CompareAndSwapInt64(&cm.poolCnt, 0, 1) {
		go func(done chan error) {
			defer atomic.StoreInt64(&cm.poolCnt, 0)
			for {
				if cm.size < cm.max {
					done = cm.grpcDial(ms, addr, done)
				} else {
					close(done)
					break
				}
			}
		}(done)
	}
}

func (cm *connPoolMeta) grpcDial(ms time.Duration, addr string, done chan error) chan error {
	//本次建连是否成功
	//失败则退出
	if !func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), ms)
		defer cancel()

		var connPool *grpc.ClientConn
		var e error
		if cm.keepaliveTime == CFGDEFKEEPALIVE {
			connPool, e = grpc.DialContext(
				ctx,
				addr,
				grpc.WithInsecure(),
				grpc.WithAuthority(addr),
				grpc.WithBlock(),
				grpc.WithReadBufferSize(cm.rbuf),
				grpc.WithWriteBufferSize(cm.wbuf),
				grpc.WithMaxMsgSize(cm.maxReceive),
			)
		} else {
			keepaliveParams := keepalive.ClientParameters{
				Time:                cm.keepaliveTime,
				Timeout:             cm.keepaliveTimeout,
				PermitWithoutStream: true}
			connPool, e = grpc.DialContext(
				ctx,
				addr,
				grpc.WithKeepaliveParams(keepaliveParams),
				grpc.WithInsecure(),
				grpc.WithAuthority(addr),
				grpc.WithBlock(),
				grpc.WithReadBufferSize(cm.rbuf),
				grpc.WithWriteBufferSize(cm.wbuf),
				grpc.WithMaxMsgSize(cm.maxReceive),
			)
		}

		if nil == e {
			if !cm.insert(&connPoolNode{cc: connPool}) {
				connPool.Close()
				done <- fmt.Errorf("can't insert conn to cm")
				return false
			} else {
				return true
			}
		} else {
			done <- e
			return false
		}
	}() {
		close(done)
	}
	return done
}

// newConnPool 创建一个新连接给
// todo: bug:链接池不满时会有链接泄露,虽然很快回收,但这给链接池初始化造成比较大开销
func (cm *connPoolMeta) newConnPool(addr string, ms time.Duration) (*grpc.ClientConn, error) {
	//先获取
	if conn, err := cm.get(addr, ms); err == nil && conn != nil {
		return conn, err
	}
	//检测池子是否填充完毕
	done := make(chan error, 1)
	if cm.size < cm.max {
		cm.fillConnPool(done, addr, ms)
	}

	for {
		select {
		case err, ok := <-done:
			{
				if ok {
					return nil, err
				}
				return cm.get(addr, ms)
			}
		default:
			{
				if conn, err := cm.get(addr, ms); err == nil && conn != nil {
					return conn, err
				}
				time.Sleep(time.Millisecond * 10)
			}
		}
	}
}

// get 从池子中获取一个链接，如果连接有问题，则剔除连接。如果连接池为空，则新建一个链接。该方法返回的连接，不一定可用
func (cm *connPoolMeta) get(addr string, ms time.Duration) (*grpc.ClientConn, error) {
	if cm.size <= 0 {
		return nil, NOUSECONN
	}
	var cc *connPoolNode
	for i := int(0); i < cm.max; i++ {
		cc = cm.next()
		if nil == cc {
			return nil, NOUSECONN
		}
		cs := cc.cc.GetState()
		// 获取的连接不可用，删除，重新获取
		//if  cs != connectivity.Idle && cs != connectivity.Connecting && cs !=connectivity.Ready{
		if cs != connectivity.Idle && cs != connectivity.Ready {
			if cs != connectivity.Connecting {
				cm.del(cc)
				cc = nil
			}
		} else {
			return cc.cc, nil
		}
	}
	//如果遍历完成,无可用链接,且都是处于Connecting,则返回当前链接
	if cc != nil && cc.cc.GetState() == connectivity.Connecting {
		return cc.cc, nil
	}

	//// 如果连接数小于最大数，重新分配连接,重试建立一个链接
	//if cm.size < cm.max {
	//	c, e := cm.newConnPool(addr, ms)
	//	if c != nil || e != nil { //只要一个不等于nil
	//		return c, e
	//	}
	//}

	return nil, NOUSECONN
}

// closeAll 关闭链表中所有节点的rpc连接
// todo  需要完善
func (cm *connPoolMeta) closeAll() {
	t := cm.head
	t.cc.Close()
	t = t.next
	for ; t != cm.head; t = t.next {
		t.cc.Close()
	}

}

//  connPool 连接池对象结构，在框架中唯一存在
type connPool struct {
	// pool addr+id 到具体连接的映射表
	pool map[string]map[string]*connPoolMeta //key:addr:key:id:value:connPool
	m    sync.RWMutex

	// o 连接选项
	o *conOption

	/*
		retry int     // 重试次数，默认2次
		timeout int  // 连接超时，单位ms
		maxmsgsize int
		max int // 最大连接数
		lc int*/
}

// newConnPool 创建connPool 管理器，需要一堆属性：
// 重试次数、连接超时等，详见ConOpt
func newConnPool(o ...connOpt) *connPool {
	c := new(connPool)

	//设置默认值
	c.o = new(conOption)
	c.o.timeout = 500
	c.o.max = 2
	c.o.lc = 120 * 1e3

	// 设置选项
	for _, opt := range o {
		opt(c.o)
	}

	// 初始化map
	c.pool = make(map[string]map[string]*connPoolMeta)
	return c
}

// get 获取连接，在接连接不稳定期，性能可能比较低 id为连接ID暂时未实现
// todo: 连接ID未真正用起来
func (c *connPool) get(addr string, id string) (*grpc.ClientConn, error) {
	//找不到对应地址的连接池
	c.m.RLock()
	s, ok := c.pool[addr]
	c.m.RUnlock()

	if !ok {
		c.m.Lock()

		_, okTmp := c.pool[addr]
		if !okTmp {
			var st connPoolMeta
			st.keepaliveTime = c.o.keepaliveTime
			st.keepaliveTimeout = c.o.keepaliveTimeout
			st.max = c.o.max
			st.wbuf = c.o.wbuf
			st.rbuf = c.o.rbuf
			st.maxReceive = c.o.maxmsgsize
			c.pool[addr] = make(map[string]*connPoolMeta)
			c.pool[addr][id] = &st
			s = c.pool[addr]
			ok = true
		} else {
			s = c.pool[addr]
			ok = true
		}
		c.m.Unlock()
	}

	//存在name 对应的addr
	c.m.RLock()
	st, ok := s[id]
	c.m.RUnlock()

	return st.alloc(addr, time.Duration(c.o.timeout)*time.Millisecond)
}

// insertAddr 在连接池中插入一个目标地址
func (c *connPool) insertAddr(addr *string, id *string, st *connPoolMeta) (bool, *connPoolMeta) {
	c.m.Lock()
	addrCon, ok := c.pool[*addr]
	if !ok {
		c.pool[*addr] = make(map[string]*connPoolMeta)
		c.pool[*addr][*id] = st

		c.m.Unlock()
		return true, nil
	} else {
		idCon, tok := addrCon[*id]
		if !tok {
			c.pool[*addr][*id] = st

			c.m.Unlock()
			return true, nil
		}
		c.m.Unlock()
		return false, idCon
	}

}

// remove 在链接池中移除一个目标地址
func (c *connPool) remove(addr *string) {
	c.m.Lock()
	tv, o := c.pool[*addr]
	if o {
		delete(c.pool, *addr)
	}

	c.m.Unlock()
	if !o {
		return
	}
	//bug: need repare
	//todo:  关闭无用连接
	if o {
		for _, mv := range tv {
			mv.closeAll()
		}

	}
}
