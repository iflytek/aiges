package stream

import (
	"github.com/xfyun/aiges/httproto/common"
	"github.com/xfyun/aiges/httproto/schemas"
	"sync"
	"time"
)

const DefaultBufferSize = 512 // read buffer size
var (
	bytePool        = NewBytePool()
	sessionPool     = NewSessionPool()
	contextPool     = sync.Pool{}
	successRespPool = sync.Pool{}
	errRespPool     = sync.Pool{}
)

func init() {
	contextPool.New = func() interface{} {
		return &schemas.Context{
			Session: nil,
			Header:  nil,
		}
	}

	successRespPool.New = func() interface{} {
		return &common.SuccessResp{
			//Message:"success",
		}
	}

	errRespPool.New = func() interface{} {
		return &common.ErrorResp{}
	}
}

//var sessionPool = NewSessionPool()

type BytePool struct {
	pool sync.Pool
}

func NewBytePool() *BytePool {
	p := &BytePool{}
	p.pool.New = func() interface{} {
		return make([]byte, DefaultBufferSize)
	}
	return p
}

func (b *BytePool) Get() []byte {
	return b.pool.Get().([]byte)
}

func (b *BytePool) Put(bf []byte) {
	b.pool.Put(bf)
}

// session pool
type SessionPoll struct {
	sync.Pool
}

func (p *SessionPoll) GetSession() *WsSession {
	return p.Get().(*WsSession)
}

func (p *SessionPoll) PutSession(session *WsSession) {
	p.Put(session)
}

func NewSessionPool() *SessionPoll {
	p := &SessionPoll{}
	p.New = func() interface{} {
		return &WsSession{}
	}
	return p
}

type ContextPool struct {
	pool sync.Pool
}

type Task func()

type TaskPool struct {
	tasks       chan func()
	size        int
	concurrency int
}

type WhenMiss func(n *Namespace, key string) (interface{}, error)
type Namespace struct {
	data     map[string]interface{}
	whenMiss WhenMiss
	timers   map[string]*time.Timer
}

func (n *Namespace) SetWithTTL(key string, val interface{}, ttl time.Duration) {
	timer := n.timers[key]
	if timer == nil {
		timer = time.AfterFunc(ttl, func() {
			delete(n.data, key)
			delete(n.timers, key)
		})
	} else {
		timer.Reset(ttl)
	}
}

func (n *Namespace) Get(key string) (interface{}, error) {
	data, ok := n.data[key]
	if !ok {
		return n.whenMiss(n, key)
	}
	return data, nil
}

func (n *Namespace) Update(key string, val interface{}) {
	n.data[key] = val
}

type Cache struct {
}

// 未命中
