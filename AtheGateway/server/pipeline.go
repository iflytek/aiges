package server

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Pipeline interface {
	//Open 代表开启管道
	Open() error
	//Push 向管道推送消息
	//isLastMsg 代表是否为最后的消息,如果为true则关闭管道
	Push(msg Message, isLastMsg bool) error
	//Close 代表关闭管道
	Close() error
	//Status 代表查看管道的状态
	Status() uint32
}

type myPipeline struct {
	nextSendMsgNo uint64             //下一个发送的消息序号
	minBreakNo    uint64             //最小的断点消息序号
	breakNoMsg    map[uint64]Message //断号的消息
	msgChan       chan Message       //消息通道
	ctx           context.Context    //上下文
	cancelFunc    context.CancelFunc //取消函数
	status        uint32             //状态
	depth         uint64             //管道深度
	rwLock        sync.RWMutex       //读写锁
	pollInterval  time.Duration      //接受的间隔时间
}

//创建一个Pipeline实例
func NewPipeline(depth uint64, pollInterval time.Duration) (Pipeline) {
	if depth <=0 {
		depth =DEFAULT_DEPTH
	}
	if depth > MAX_DEPTH {
		//return nil, fmt.Errorf("depth (%+v) paramter > max path (%+v)", depth, MAX_DEPTH)
		depth = MAX_DEPTH
	}
	if pollInterval == 0 {
		pollInterval = time.Microsecond
	}
	pip := &myPipeline{
		breakNoMsg:   make(map[uint64]Message, depth),
		status:       STATUS_ORIGINAL,
		depth:        depth,
		msgChan:      make(chan Message, depth),
		pollInterval: pollInterval,
	}
	return pip
}

func (pip *myPipeline) Open() error {
	pip.rwLock.Lock()
	defer pip.rwLock.Unlock()
	err := checkStatus(pip.status, STATUS_STARTING, nil)
	if err != nil {
		return err
	}
	pip.ctx, pip.cancelFunc = context.WithCancel(context.Background())
	atomic.StoreUint32(&pip.status, STATUS_STARTED)
	go func() {
		for {
			select {
			case <-pip.ctx.Done():
				pip.prepareToClose()
				goto Loop
			case msg := <-pip.msgChan:
				pip.sendMsg(msg)
			//case <-time.After(pip.pollInterval):
			//default:
			}
		}
	Loop:
		pip.clearBreakNoMsg()
	}()
	return nil
}

//restMinBreakNo 代表重置断号
func (pip *myPipeline) resetMinBreakNo() {
	isFist := true
	for k, _ := range pip.breakNoMsg {
		if isFist {
			pip.minBreakNo = k
			isFist = false
			continue
		}
		if pip.minBreakNo > k {
			pip.minBreakNo = k
		}
	}
}

//sendMsg  代表发送消息
func (pip *myPipeline) sendMsg(msg Message) {
	if msg.MsgNo() == pip.nextSendMsgNo {
		pip.doSendMsg(msg)
		return
	}
	if msg.MsgNo() != pip.nextSendMsgNo {
		if pip.nextSendMsgNo >= msg.MsgNo() {
			//过期消息已经被丢弃
			return
		}
		pip.breakNoMsg[msg.MsgNo()] = msg
		pip.resetMinBreakNo()
		if len(pip.breakNoMsg) == int(pip.depth) {
			pip.doSendMsg(pip.breakNoMsg[pip.minBreakNo])
			if len(pip.breakNoMsg) > 0 {
				pip.resetMinBreakNo()
			}
		}
	}
}

//doSendMsg  代表发送消息
func (pip *myPipeline) doSendMsg(msg Message) error {
	msg.Send()
	if len(pip.breakNoMsg) > 0 {
		delete(pip.breakNoMsg, msg.MsgNo())
	}
	nextMsgNo := msg.MsgNo() + 1
	atomic.SwapUint64(&pip.nextSendMsgNo, nextMsgNo)
	if msg, exist := pip.breakNoMsg[nextMsgNo]; len(pip.breakNoMsg) > 0 && exist {
		pip.doSendMsg(msg)
	}
	return nil
}

//prepareToStop 代表用于为停止管道做准备
func (pip *myPipeline) prepareToClose() {
	pip.clearMsgChan()
	pip.rwLock.Lock()
	defer pip.rwLock.Unlock()
	if atomic.CompareAndSwapUint32(&pip.status, STATUS_STOPPING, STATUS_STOPPED) {
		close(pip.msgChan)
	}
}

//clearMsgChan 代表清空消息通道
func (pip *myPipeline) clearMsgChan() {
	for ; len(pip.msgChan) > 0; {
		pip.sendMsg(<-pip.msgChan)
	}
}

//clearBreakNoMsg 代表清空断点的消息
func (pip *myPipeline) clearBreakNoMsg() {
	if len(pip.breakNoMsg) == 0 {
		return
	}
	for i := int(atomic.LoadUint64(&pip.nextSendMsgNo)); len(pip.breakNoMsg) > 0; i = i + 1 {
		if msg, exist := pip.breakNoMsg[uint64(i)]; exist {
			msg.Send()
			delete(pip.breakNoMsg, uint64(i))
		}
	}
}

func (pip *myPipeline) Push(msg Message, isLastMsg bool) error {
	if msg == nil {
		return errors.New("msg is nil")
	}
	pip.rwLock.RLock()
	defer pip.rwLock.RUnlock()
	err := checkStatus(pip.status, STATUS_STOPPING, nil)
	if err != nil {
		return err
	}
	pip.msgChan <- msg
	if isLastMsg {
		atomic.StoreUint32(&pip.status, STATUS_STOPPING)
		pip.cancelFunc()
	}
	return nil
}

func (pip *myPipeline) Close() error {
	pip.rwLock.Lock()
	defer pip.rwLock.Unlock()
	err := checkStatus(pip.status, STATUS_STOPPING, nil)
	if err == nil {
		pip.cancelFunc()
		atomic.StoreUint32(&pip.status, STATUS_STOPPING)
		return nil
	}
	return err
}

func (pip *myPipeline) Status() uint32 {
	return atomic.LoadUint32(&pip.status)
}
