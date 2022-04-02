package buffer

import (
	"context"
	"github.com/xfyun/aiges/frame"
	"sync"
	"time"
)

/*
	数据缓存及排序模块;用于引擎服务对上下行数据的缓存或排序操作;
	1.分块数据的排序
	2.乱序数据的超时获取(阻塞&非阻塞)
	3.多协程并发写,单协程有序读
	4.缓冲区状态查询
	5.可获取合并数据
	6.异常场景中断读超时等待(功能优化):
	7.多数据流输入输出;
*/

type metaStream map[uint] /*frameId*/ DataMeta
type MultiBuf struct {
	order     chan DataMeta // 有序缓冲区,待读取
	disorder  map[string] /*DataId*/ metaStream
	nextFrame map[string] /*DataId*/ uint /*current index*/
	mutex     sync.RWMutex                // 缓冲区数据读写锁;

	orderSize  uint      // 有序缓冲区大小
	baseFrame  uint      // 数据流首帧起始id
	baseCache  uint      // baseFrame备份
	waitTimeMs uint      // 排序读取超时等待时间
	signal     chan bool // 宿主存活状态同步信号;
}

func (mb *MultiBuf) Init(time uint, size uint, base uint) {
	mb.orderSize = size
	mb.baseFrame = base
	mb.baseCache = base
	mb.waitTimeMs = time
	mb.order = make(chan DataMeta, mb.orderSize)
	mb.disorder = make(map[string]metaStream)
	mb.nextFrame = make(map[string]uint)
	mb.signal = make(chan bool, 1) // chan大小为1,防止未buf.waitTime场景 协程退出时写chan阻塞;
	return
}

func (mb *MultiBuf) Fini() {
	close(mb.order)
	close(mb.signal)
	for id := range mb.disorder {
		delete(mb.disorder, id)
	}
	for id := range mb.nextFrame {
		delete(mb.nextFrame, id)
	}
	return
}

// call set before write and read;
func (mb *MultiBuf) SetBase(base uint) {
	mb.baseFrame = base
}

func (mb *MultiBuf) WriteData(input []DataMeta) (errNum int, err error) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	for _, v := range input {
		// 不存在,首次写入
		curIndex, exist := mb.nextFrame[v.DataId]
		if !exist {
			if v.FrameId == mb.baseFrame || v.Status == DataStatusFirst || v.Status == DataStatusOnce {
				// 首帧数据
				mb.order <- v
				mb.nextFrame[v.DataId] = v.FrameId + 1
			} else {
				stream := make(map[uint]DataMeta)
				stream[v.FrameId] = v
				mb.disorder[v.DataId] = stream
				mb.nextFrame[v.DataId] = mb.baseFrame
			}
		}
		// 已存在,非首次写入
		if exist && v.FrameId >= curIndex {
			if curIndex == v.FrameId {
				select {
				case mb.order <- v:
					mb.nextFrame[v.DataId] = v.FrameId + 1
					// check disorder, move disorder to order Chan
					if stream, exist := mb.disorder[v.DataId]; exist {
						next := mb.nextFrame[v.DataId]
					loopBreak:
						for {
							if meta, exist := stream[next]; exist {
								select {
								case mb.order <- meta:
									delete(stream, next)
									next++
									mb.nextFrame[v.DataId] = next
								default:
									break loopBreak
								}
								continue
							}
							break
						}
					}
					continue
				default:
					// 有序chan满
				}
			}
			// 写入无序缓冲区
			if _, exist := mb.disorder[v.DataId]; !exist {
				mb.disorder[v.DataId] = make(map[uint]DataMeta)
			}
			mb.disorder[v.DataId][v.FrameId] = v
		}
	}

	return
}

// ReadDataNonBlock do not merge data & nonblock,
// just read current minimal frame data and move the next frame flag
func (mb *MultiBuf) ReadDataNonBlock() (output []DataMeta, errNum int, errInfo error) {
	mb.migrate()
	// nonBlock调用;
	select {
	case meta, open := <-mb.order:
		if open {
			output = make([]DataMeta, 0, 1)
			output = append(output, meta)
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
func (mb *MultiBuf) ReadDataWithTime(timeout uint /*ms*/) (output []DataMeta, errNum int, errInfo error) {
	// 校验disorder是否存在可迁移数据;
	mb.migrate()
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
		defer cancel()
		select {
		case <-mb.signal:
			errNum = frame.AigesErrorSeqChanClosed
			errInfo = frame.ErrorSeqChanClosed
			return
		case meta, open := <-mb.order:
			if open {
				output = make([]DataMeta, 0, 1)
				output = append(output, meta)
			} else {
				errNum = frame.AigesErrorSeqChanClosed
				errInfo = frame.ErrorSeqChanClosed
			}
			return
		case <-ctx.Done():
			// return if order channel exist data ,else pass to call func ReadDataNonBlock
		}
	}

	output, errNum, errInfo = mb.ReadDataNonBlock()
	return
}

func (mb *MultiBuf) ReadMergeData(merge bool, timeout *uint /*ms*/) (output []DataMeta, errNum int, errInfo error) {
	mb.migrate()
	tmpData := make(map[string]DataMeta)
	// select for timeout wait, block read
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
	case <-ctx.Done():
		mb.degrade()
	case meta, open := <-mb.order:
		if open {
			if !merge {
				output = append(output, meta)
				return
			}
			tmpData[meta.DataId] = meta
		} else {
			errNum = frame.AigesErrorSeqChanClosed
			errInfo = frame.ErrorSeqChanClosed
			return
		}
	}

	// select for merge order data
	for {
		select {
		case meta, open := <-mb.order:
			if open {
				if !merge {
					output = append(output, meta)
					return
				}
				stream, ok := tmpData[meta.DataId]
				if ok && stream.Data != nil { // merge the same stream's data
					if meta.Data != nil {
						stream.Data = append(stream.Data.([]byte), (meta.Data).([]byte)[0:]...)
					}
					stream.FrameId = meta.FrameId
					stream.Status = meta.Status
					tmpData[meta.DataId] = stream
				} else {
					tmpData[meta.DataId] = meta
				}
			} else {
				errNum = frame.AigesErrorSeqChanClosed
				errInfo = frame.ErrorSeqChanClosed
				return
			}
		default:
			// order channel empty,  return merged data
			if len(tmpData) > 0 {
				output = make([]DataMeta, 0, len(tmpData))
				for _, data := range tmpData {
					output = append(output, data)
				}
			} else {
				errNum = frame.AigesErrorBufferEmpty
				errInfo = frame.ErrorSeqBufferEmpty
			}
			return
		}
	}
}

func (mb *MultiBuf) Release() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	close(mb.signal)
	close(mb.order)
	mb.signal = make(chan bool, 1)
	mb.order = make(chan DataMeta, mb.orderSize)
	mb.baseFrame = mb.baseCache
	for id := range mb.disorder {
		delete(mb.disorder, id)
	}
	for id := range mb.nextFrame {
		delete(mb.nextFrame, id)
	}
}

func (mb *MultiBuf) Signal() {
	mb.signal <- true
}

// 降级迁移,丢弃超时不至的数据分片
func (mb *MultiBuf) degrade() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	if len(mb.order) == 0 && len(mb.disorder) > 0 {
		// 对所有stream进行降级迁移
		for id, stream := range mb.disorder {
			if frame := minimalFrame(stream); frame >= 0 {
				// find stream's minimal existed frameId
				mb.nextFrame[id] = uint(frame)
			loopBreak:
				for {
					if meta, exist := mb.disorder[id][mb.nextFrame[id]]; exist {
						select {
						case mb.order <- meta:
							delete(mb.disorder[id], mb.nextFrame[id])
							mb.nextFrame[id]++
						default:
							break loopBreak
						}
						continue
					}
					break
				}
			}
		}
	}
	return
}

func (mb *MultiBuf) degradeRead() (output []DataMeta, errNum int, errInfo error) {
	mb.degrade()
	// select read order channel data
	select {
	case meta, open := <-mb.order:
		if open {
			output = make([]DataMeta, 0, 1)
			output = append(output, meta)
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

//func (mb *MultiBuf) degradeMerge() (output []DataMeta, errNum int, errInfo error) {
//	// 降级合并读
//	mb.degrade()
//	return
//}

// 迁移disorder中有序数据至order channel
func (mb *MultiBuf) migrate() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	for streamId := range mb.nextFrame {
		if stream, exist := mb.disorder[streamId]; exist {
			// data stream disorder meta move to order channel
		loopBreak:
			for {
				if meta, exist := stream[mb.nextFrame[streamId]]; exist {
					select {
					case mb.order <- meta:
						delete(stream, mb.nextFrame[streamId])
						mb.nextFrame[streamId]++
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

// return -1 if stream is nil or empty,
// else return minimal frameId of metaStream
func minimalFrame(stream metaStream) (frameId int) {
	frameId = -1
	for i, _ := range stream {
		if frameId < 0 || frameId > int(i) {
			frameId = int(i)
		}
	}
	return
}
