package buffer

import (
	"context"
	"github.com/xfyun/aiges/frame"
	"sync"
	"time"
)

// 业务包排序逻辑：取消底层数据流排序及合并操作;
//type dataStream map[uint] /*frameId*/ DataMeta
type DataBuf struct {
	order      chan []DataMeta // 有序缓冲区,待读取
	disorder   map[uint] /*seqId*/ []DataMeta
	nextFrame  uint         // 当前有序帧
	mutex      sync.RWMutex // 缓冲区数据读写锁;
	baseCache  uint         // baseFrame备份
	orderSize  uint         // 有序缓冲区大小
	waitTimeMs uint         // 排序读取超时等待时间
	signal     chan bool    // 宿主存活状态同步信号;
}

func (mb *DataBuf) Init(time uint, size uint, base uint) {
	mb.orderSize = size
	mb.baseCache = base
	mb.nextFrame = base
	mb.waitTimeMs = time
	mb.order = make(chan []DataMeta, mb.orderSize)
	mb.disorder = make(map[uint][]DataMeta)
	mb.signal = make(chan bool, 1) // chan大小为1,防止未buf.waitTime场景 协程退出时写chan阻塞;
	return
}

func (mb *DataBuf) Fini() {
	close(mb.order)
	close(mb.signal)
	for id := range mb.disorder {
		delete(mb.disorder, id)
	}
	return
}

// call set before write and read;
func (mb *DataBuf) SetBase(base uint) {
	mb.nextFrame = base
}

func (mb *DataBuf) WriteData(seq uint, input []DataMeta) (errNum int, err error) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	// data -> 无序缓冲区 -> 有序缓冲区
	if seq < mb.nextFrame {
		return // drop 无效数据
	}
	mb.disorder[seq] = input
loopBreak:
	for {
		if metas, exist := mb.disorder[mb.nextFrame]; exist {
			if metas == nil || len(metas) == 0 {
				delete(mb.disorder, mb.nextFrame)
				mb.nextFrame++ // 防止消费端读取nil,影响超时wait判定逻辑
				continue
			}
			select {
			case mb.order <- metas:
				delete(mb.disorder, mb.nextFrame)
				mb.nextFrame++
			default:
				break loopBreak
			}
			continue
		}
		break
	}

	return
}

// ReadDataNonBlock do not merge data & nonblock,
// just read current minimal frame data and move the next frame flag
func (mb *DataBuf) ReadDataNonBlock() (output []DataMeta, errNum int, errInfo error) {
	mb.migrate()
	// nonBlock调用;
	select {
	case meta, open := <-mb.order:
		if open {
			output = meta
		} else {
			errNum = frame.AigesErrorSeqChanClosed
			errInfo = frame.ErrorSeqChanClosed
		}
	default:
		// 获取非有序缓冲区数据&同步有序缓冲区;
		output, errNum, errInfo = mb.degradeRead()
	}
	return
}

// do not merge data
func (mb *DataBuf) ReadDataWithTime(timeout *uint /*ms*/) (output []DataMeta, errNum int, errInfo error) {
	// 校验disorder是否存在可迁移数据;
	mb.migrate()

	// default wait time
	tm := mb.waitTimeMs
	if timeout != nil {
		tm = *timeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tm)*time.Millisecond)
	defer cancel()
	select {
	case <-mb.signal:
		errNum = frame.AigesErrorSeqChanClosed
		errInfo = frame.ErrorSeqChanClosed
		return
	case meta, open := <-mb.order:
		if open {
			output = meta
		} else {
			errNum = frame.AigesErrorSeqChanClosed
			errInfo = frame.ErrorSeqChanClosed
		}
		return
	case <-ctx.Done():
		// return if order channel exist data ,else pass to call func ReadDataNonBlock
	}

	output, errNum, errInfo = mb.ReadDataNonBlock()
	return
}

func (mb *DataBuf) Release() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	close(mb.signal)
	close(mb.order)
	mb.signal = make(chan bool, 1)
	mb.order = make(chan []DataMeta, mb.orderSize)
	mb.nextFrame = mb.baseCache
	for id := range mb.disorder {
		delete(mb.disorder, id)
	}
}

func (mb *DataBuf) Signal() {
	select {
	case mb.signal <- true:
	default:
		return
	}
}

// 降级迁移,丢弃超时不至的数据分片
func (mb *DataBuf) degrade() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	if len(mb.order) == 0 && len(mb.disorder) > 0 {
		// 对所有stream进行降级迁移
		if frame := minFrame(mb.disorder); frame >= 0 {
			mb.nextFrame = uint(frame)
		loopBreak:
			for {
				if meta, exist := mb.disorder[mb.nextFrame]; exist {
					if meta == nil || len(meta) == 0 {
						delete(mb.disorder, mb.nextFrame)
						mb.nextFrame++ // 防止消费端读取nil,影响超时wait判定逻辑
						continue
					}
					select {
					case mb.order <- meta:
						delete(mb.disorder, mb.nextFrame)
						mb.nextFrame++
					default:
						break loopBreak
					}
					continue
				}
				break
			}
		}
	}
	return
}

func (mb *DataBuf) degradeRead() (output []DataMeta, errNum int, errInfo error) {
	mb.degrade()
	// select read order channel data
	select {
	case meta, open := <-mb.order:
		if open {
			output = meta
		} else {
			errNum = frame.AigesErrorSeqChanClosed
			errInfo = frame.ErrorSeqChanClosed
		}
	default:
		// nothing to do, empty.
		errNum = frame.AigesErrorBufferEmpty
		errInfo = frame.ErrorSeqBufferEmpty
	}
	return
}

// 迁移disorder中有序数据至order channel
func (mb *DataBuf) migrate() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	// data stream disorder meta move to order channel
loopBreak:
	for {
		if meta, exist := mb.disorder[mb.nextFrame]; exist {
			if meta == nil || len(meta) == 0 {
				delete(mb.disorder, mb.nextFrame)
				mb.nextFrame++ // 防止消费端读取nil,影响超时wait判定逻辑
				continue
			}
			select {
			case mb.order <- meta:
				delete(mb.disorder, mb.nextFrame)
				mb.nextFrame++
			default:
				break loopBreak
			}
			continue
		}
		break
	}
	return
}

// return -1 if stream is nil or empty,
// else return minimal frameId of metaStream
func minFrame(stream map[uint][]DataMeta) (frameId int) {
	frameId = -1
	for i, _ := range stream {
		if frameId < 0 || frameId > int(i) {
			frameId = int(i)
		}
	}
	return
}
