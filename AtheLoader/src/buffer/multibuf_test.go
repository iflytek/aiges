package buffer

import (
	"bytes"
	"frame"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

/*
	测试场景
	1. 单块音频 & baseId & status验证
	2. 两块音频(乱序)
	3. 多块音频(乱序)
	4. 并发写入(乱序)
	5. 超时验证
	6. 数据一致性
	7. 多个数据流场景
	8. 异常超时中断
*/

const (
	defaultTimeoutDeviation = 10   // 超时误差;ms
	defaultTimeout          = 1000 // 缺省超时时间
	defaultSeqBufLen        = 100  // 缺省有序缓冲区长度
	defaultBaseId           = 0    // 缺省基准Id
	defaultStreamLen        = 10   // 测试数据流长度
)

// 缺省定义数据流
type TestMeta struct {
	id     uint
	data   []byte
	status DataStatus
}

var TestDataStream []TestMeta
var TestDataStream2 []TestMeta // 数据实体区分于TestDataStream,防止数据流数据异常无法验证
var TestData []byte
var TestData2 []byte

func init() {
	TestData = make([]byte, 0, defaultStreamLen)
	TestData2 = make([]byte, 0, defaultStreamLen)
	TestDataStream = make([]TestMeta, defaultStreamLen)
	TestDataStream2 = make([]TestMeta, defaultStreamLen)
	for index := 0; index < defaultStreamLen; index++ {
		// 数据流1
		TestDataStream[index].id = uint(index)
		TestDataStream[index].data = make([]byte, 0, 1)
		TestDataStream[index].data = append(TestDataStream[index].data, byte(index))
		TestDataStream[index].status = DataStatusContinue
		TestData = append(TestData, TestDataStream[index].data...)
		if index == 0 {
			TestDataStream[index].status = DataStatusFirst
		}
		if index == defaultStreamLen-1 {
			TestDataStream[index].status = DataStatusLast
		}

		// 数据流2
		TestDataStream2[index].id = uint(index)
		TestDataStream2[index].data = make([]byte, 0, 1)
		TestDataStream2[index].data = append(TestDataStream2[index].data, byte(index+1))
		TestDataStream2[index].status = DataStatusContinue
		TestData2 = append(TestData2, TestDataStream2[index].data...)
		if index == 0 {
			TestDataStream2[index].status = DataStatusFirst
		}
		if index == defaultStreamLen-1 {
			TestDataStream2[index].status = DataStatusLast
		}
	}
	return
}

// 单个数据写入
func TestSingleData(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)

	// Example: 单数据基准id写入
	// 期望：即时读取
	buf.WriteData([]DataMeta{{nil, "0", 0, DataStatusLast, DataText, "", "", nil}})
	rdTime := time.Now()
	outPut, _, err := buf.ReadDataWithTime(defaultTimeout)
	rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("基准id,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("基准id,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != 0 || outPut[0].Status != DataStatusLast {
		t.Errorf("基准id,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}
	buf.Release()

	// Example: 单数据非基准id异常写入
	// 期望：超时读取
	buf.WriteData([]DataMeta{{nil, "0", 5, DataStatusLast, DataText, "", "", nil}})
	rdTime = time.Now()
	outPut, _, err = buf.ReadDataWithTime(defaultTimeout)
	rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if math.Abs(float64(rdPerf-defaultTimeout)) > defaultTimeoutDeviation {
		t.Errorf("非基准id,读耗时异常未超时等待: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("非基准id,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != 5 || outPut[0].Status != DataStatusLast {
		t.Errorf("非基准id,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}
	buf.Release()

	// Example: 修改baseId读写
	// 期望：即时读取
	buf.SetBase(5)
	buf.WriteData([]DataMeta{{nil, "0", 5, DataStatusLast, DataText, "", "", nil}})
	rdTime = time.Now()
	outPut, _, err = buf.ReadDataWithTime(defaultTimeout)
	rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("非基准id,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("非基准id,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != 5 || outPut[0].Status != DataStatusLast {
		t.Errorf("非基准id,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}
	buf.Release()

	// Example: 重复写入
	// 期望：即时正常读
	buf.WriteData([]DataMeta{{nil, "0", 0, DataStatusLast, DataText, "", "", nil}})
	buf.WriteData([]DataMeta{{nil, "0", 0, DataStatusLast, DataText, "", "", nil}})
	rdTime = time.Now()
	outPut, _, err = buf.ReadDataWithTime(defaultTimeout)
	rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("重复写入读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("重复写入读,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != 0 || outPut[0].Status != DataStatusLast {
		t.Errorf("重复写入读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}

	// Example: 持续读操作
	// 期望：报错缓冲区空，TODO 不等待超时
	outPut, _, err = buf.ReadDataWithTime(defaultTimeout)
	if err != frame.ErrorSeqBufferEmpty {
		t.Errorf("异常读,错误不符预期")
	}
	buf.Release()
	buf.Fini()
	return
}

// 两个数据写入(乱序)
func TestDoubleData(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)

	// Example: meta读
	buf.WriteData([]DataMeta{{nil, "0", 1, DataStatusLast, DataText, "", "", nil}})
	buf.WriteData([]DataMeta{{nil, "0", 0, DataStatusFirst, DataText, "", "", nil}})
	rdTime := time.Now()
	outPut, _, err := buf.ReadDataWithTime(defaultTimeout)
	rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("乱序写入读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("乱序写入读,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != 0 || outPut[0].Status != DataStatusFirst {
		t.Errorf("乱序写入读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}

	rdTime = time.Now()
	outPut, _, err = buf.ReadDataWithTime(defaultTimeout)
	rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("乱序写入读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("乱序写入读,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != 1 || outPut[0].Status != DataStatusLast {
		t.Errorf("乱序写入读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}
	buf.Release()

	// Example: 合并读
	buf.WriteData([]DataMeta{{nil, "0", 1, DataStatusLast, DataText, "", "", nil}})
	buf.WriteData([]DataMeta{{nil, "0", 0, DataStatusFirst, DataText, "", "", nil}})
	rdTime = time.Now()
	outPut, _, err = buf.ReadMergeData()
	rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("乱序写入合并读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("乱序写入合并读,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != 1 || outPut[0].Status != DataStatusLast {
		t.Errorf("乱序写入合并读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}
	buf.Release()
	buf.Fini()
	return
}

// 多个数据写入(乱序) & 数据一致性
func TestMultiData(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)

	// Example: 乱序写&合并读
	streamTmp := make([]TestMeta, defaultStreamLen)
	copy(streamTmp, TestDataStream)
	streamCnt := defaultStreamLen
	for ; streamCnt > 0; streamCnt-- {
		randId := rand.Intn(streamCnt)
		buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
			streamTmp[randId].status, DataText, "", "", nil}})
		streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
	}

	// 合并读
	rdTime := time.Now()
	outPut, _, err := buf.ReadMergeData()
	rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("乱序写入合并读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("乱序写入合并读,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != defaultStreamLen-1 || outPut[0].Status != DataStatusLast {
		t.Errorf("乱序写入合并读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}
	// 数据一致性校验
	if bytes.Compare(TestData, outPut[0].Data.([]byte)) != 0 {
		t.Fatal(TestData, outPut[0].Data.([]byte))
		t.Errorf("乱序写入合并读,数据一致性异常")
	}
	buf.Release()

	// Example: 乱序写&分段读
	streamTmp = make([]TestMeta, defaultStreamLen)
	copy(streamTmp, TestDataStream)
	streamCnt = defaultStreamLen
	for ; streamCnt > 0; streamCnt-- {
		randId := rand.Intn(streamCnt)
		buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
			streamTmp[randId].status, DataText, "", "", nil}})
		streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
	}
	// meta读
	index := 0
	status := DataStatusContinue
	for ; status != DataStatusLast; index++ {
		rdTime = time.Now()
		outPut, _, err = buf.ReadDataWithTime(defaultTimeout)
		rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
		if rdPerf > defaultTimeoutDeviation {
			t.Errorf("乱序写入meta读,读耗时过长: %d ms", rdPerf)
		} else if err != nil {
			t.Errorf("乱序写入meta读,读操作异常: %s", err.Error())
		} else if outPut[0].FrameId != uint(index) {
			t.Errorf("乱序写入meta读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
		} else if bytes.Compare(outPut[0].Data.([]byte), TestDataStream[index].data) != 0 { // 数据一致性校验;
			t.Errorf("乱序写入meta读,数据一致性异常")
		}
		status = outPut[0].Status
	}
	buf.Release()
	buf.Fini()
	return
}

// 并发写入
func TestConcurrence(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)
	wg := sync.WaitGroup{}

	// Example: 并发乱序写&合并读
	for rtId := 0; rtId < 10; rtId++ {
		wg.Add(1)
		go func() {
			// Example: 并发乱序写
			streamTmp := make([]TestMeta, defaultStreamLen)
			copy(streamTmp, TestDataStream)
			streamCnt := defaultStreamLen
			for ; streamCnt > 0; streamCnt-- {
				randId := rand.Intn(streamCnt)
				buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
					streamTmp[randId].status, DataText, "", "", nil}})
				streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	rdTime := time.Now()
	outPut, _, err := buf.ReadMergeData()
	rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("并发乱序写入合并读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("并发乱序写入合并读,读操作异常: %s", err.Error())
	} else if outPut[0].FrameId != defaultStreamLen-1 || outPut[0].Status != DataStatusLast {
		t.Errorf("并发乱序写入合并读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
	}
	// 数据一致性校验
	if bytes.Compare(TestData, outPut[0].Data.([]byte)) != 0 {
		t.Errorf("乱序写入合并读,数据一致性异常")
	}
	buf.Release()

	// Example: 并发乱序写&重复写&分段读
	for rtId := 0; rtId < 10; rtId++ {
		wg.Add(1)
		go func() {
			// Example: 并发乱序写
			streamTmp := make([]TestMeta, defaultStreamLen)
			copy(streamTmp, TestDataStream)
			streamCnt := defaultStreamLen
			for ; streamCnt > 0; streamCnt-- {
				randId := rand.Intn(streamCnt)
				buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
					streamTmp[randId].status, DataText, "", "", nil}})
				streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	index := 0
	status := DataStatusContinue
	for ; status != DataStatusLast; index++ {
		rdTime = time.Now()
		outPut, _, err = buf.ReadDataWithTime(defaultTimeout)
		rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
		if rdPerf > defaultTimeoutDeviation {
			t.Errorf("乱序写入meta读,读耗时过长: index:%d, %d ms, %v", index, rdPerf, outPut)
			break
		} else if err != nil {
			t.Errorf("乱序写入meta读,读操作异常: %s", err.Error())
		} else if outPut[0].FrameId != uint(index) {
			t.Errorf("乱序写入meta读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, status)
		} else if bytes.Compare(outPut[0].Data.([]byte), TestDataStream[index].data) != 0 { // 数据一致性校验;
			t.Errorf("乱序写入meta读,数据一致性异常")
		}
		status = outPut[0].Status
	}
	buf.Release()
	buf.Fini()
	return
}

// 乱序超时
func TestReadTimeout(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)
	var indexExpect uint

	// Example: 乱序写&合并读超时
	streamTmp := make([]TestMeta, defaultStreamLen)
	copy(streamTmp, TestDataStream)
	for streamCnt := defaultStreamLen; streamCnt > 0; streamCnt-- {
		randId := rand.Intn(streamCnt)
		buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
			streamTmp[randId].status, DataText, "", "", nil}})

		// 合并读校验
		rdTime := time.Now()
		outPut, _, err := buf.ReadMergeData()
		rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)

		// 顺序读场景
		if streamTmp[randId].id == indexExpect {
			if rdPerf > defaultTimeoutDeviation {
				t.Errorf("乱序超时合并读,读耗时过长: %d ms", rdPerf)
			} else if err != nil {
				t.Errorf("乱序超时合并读,读操作异常: %s", err.Error())
			} else if outPut[0].FrameId != indexExpect || outPut[0].Status != streamTmp[randId].status {
				t.Errorf("乱序超时合并读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
			}
			indexExpect++
		} else if streamTmp[randId].id > indexExpect { // 超时读有数据
			if math.Abs(float64(rdPerf-defaultTimeout)) > defaultTimeoutDeviation {
				t.Errorf("乱序超时合并读,读耗时异常未超时等待: %d ms", rdPerf)
			} else if err != nil {
				t.Errorf("乱序超时合并读,读操作异常: %s", err.Error())
			} else if outPut[0].FrameId != streamTmp[randId].id || outPut[0].Status != streamTmp[randId].status {
				t.Errorf("乱序超时合并读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
			}
			indexExpect = outPut[0].FrameId + 1
		} else { // 数据写丢弃,超时无数据
			if math.Abs(float64(rdPerf-defaultTimeout)) > defaultTimeoutDeviation {
				t.Errorf("乱序超时合并读,读耗时异常未超时等待: %d ms", rdPerf)
			} else if err != frame.ErrorSeqBufferEmpty {
				t.Errorf("乱序超时合并读,读操作异常未报错缓冲区空")
			}
		}
		streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
	}
	buf.Release()

	// Example: 乱序写&分段读超时
	indexExpect = 0
	streamTmp = make([]TestMeta, defaultStreamLen)
	copy(streamTmp, TestDataStream)
	for streamCnt := defaultStreamLen; streamCnt > 0; streamCnt-- {
		randId := rand.Intn(streamCnt)
		buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
			streamTmp[randId].status, DataText, "", "", nil}})

		// 分段读校验
		rdTime := time.Now()
		outPut, _, err := buf.ReadDataWithTime(defaultTimeout)
		rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)

		// 顺序读场景
		if streamTmp[randId].id == indexExpect {
			if rdPerf > defaultTimeoutDeviation {
				t.Errorf("乱序超时分段读,读耗时过长: %d ms", rdPerf)
			} else if err != nil {
				t.Errorf("乱序超时分段读,读操作异常: %s", err.Error())
			} else if outPut[0].FrameId != indexExpect || outPut[0].Status != streamTmp[randId].status {
				t.Errorf("乱序超时分段读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
			}
			indexExpect++
		} else if streamTmp[randId].id > indexExpect { // 超时读有数据
			if math.Abs(float64(rdPerf-defaultTimeout)) > defaultTimeoutDeviation {
				t.Errorf("乱序超时分段读,读耗时异常未超时等待: %d ms", rdPerf)
			} else if err != nil {
				t.Errorf("乱序超时分段读,读操作异常: %s", err.Error())
			} else if outPut[0].FrameId != streamTmp[randId].id || outPut[0].Status != streamTmp[randId].status {
				t.Errorf("乱序超时合并读,读数据异常 id:%d, 状态:%d", outPut[0].FrameId, outPut[0].Status)
			}
			indexExpect = outPut[0].FrameId + 1
		} else { // 数据写丢弃,超时无数据
			if math.Abs(float64(rdPerf-defaultTimeout)) > defaultTimeoutDeviation {
				t.Errorf("乱序超时分段读,读耗时异常未超时等待: %d ms", rdPerf)
			} else if err != frame.ErrorSeqBufferEmpty {
				t.Errorf("乱序超时分段读,读操作异常未报错缓冲区空")
			}
		}
		streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
	}
	buf.Release()
	buf.Fini()
	return
}

func TestEmptyReadTimeout(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)
	var timeoutMs uint = 4000
	rdTime := time.Now()
	_, _, _ = buf.ReadDataWithTime(timeoutMs)
	rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if math.Abs(float64(rdPerf)-float64(timeoutMs)) > defaultTimeoutDeviation {
		t.Errorf("无数据写入等待超时场景: %d ms", rdPerf)
		return
	}
	return
}

// 读取超时中断
func TestRdInterrupt(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)
	var timeoutMs uint = 4000
	buf.Signal()
	rdTime := time.Now()
	_, _, _ = buf.ReadDataWithTime(timeoutMs)
	rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("中断超时等待场景: %d ms", rdPerf)
		return
	}
	return
}

func TestMultiStream(t *testing.T) {
	buf := MultiBuf{}
	buf.Init(defaultTimeout, defaultSeqBufLen, defaultBaseId)
	wg := sync.WaitGroup{}

	// 数据流"0"写入
	for rtId := 0; rtId < 10; rtId++ {
		wg.Add(1)
		go func() {
			// Example: 并发乱序写
			streamTmp := make([]TestMeta, defaultStreamLen)
			copy(streamTmp, TestDataStream)
			streamCnt := defaultStreamLen
			for ; streamCnt > 0; streamCnt-- {
				randId := rand.Intn(streamCnt)
				buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
					streamTmp[randId].status, DataText, "", "", nil}})
				streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
			}
			wg.Done()
		}()
	}

	// 数据流"1"写入
	for rtId := 0; rtId < 10; rtId++ {
		wg.Add(1)
		go func() {
			// Example: 并发乱序写
			streamTmp := make([]TestMeta, defaultStreamLen)
			copy(streamTmp, TestDataStream2)
			streamCnt := defaultStreamLen
			for ; streamCnt > 0; streamCnt-- {
				randId := rand.Intn(streamCnt)
				buf.WriteData([]DataMeta{{streamTmp[randId].data, "1", streamTmp[randId].id,
					streamTmp[randId].status, DataText, "", "", nil}})
				streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// 读取 & 数据流完成性校验
	rdTime := time.Now()
	outPut, _, err := buf.ReadMergeData()
	rdPerf := time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("多数据流并发乱序写入合并读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("多数据流并发乱序写入合并读,读操作异常: %s", err.Error())
	}

	// 数据一致性校验
	eq, eq2 := false, false
	for _, stream := range outPut {
		if stream.DataId == "0" && bytes.Compare(TestData, stream.Data.([]byte)) == 0 &&
			stream.FrameId == defaultStreamLen-1 && stream.Status == DataStatusLast {
			eq = true
		} else if stream.DataId == "1" && bytes.Compare(TestData2, stream.Data.([]byte)) == 0 &&
			stream.FrameId == defaultStreamLen-1 && stream.Status == DataStatusLast {
			eq2 = true
		}
	}
	if !eq || !eq2 {
		t.Errorf("多数据流数据输入输出异常, outPut: %v", outPut)
	}
	buf.Release()

	/*场景：
	输入：
		stream1 数据流正常完成写入
		stream2 数据流异常未完成写入
	期望输出：
		stream1 未超时即时读取
		stream2 超时降级排序读取
	*/
	// 数据流"0"写入
	for rtId := 0; rtId < 10; rtId++ {
		wg.Add(1)
		go func() {
			// Example: 并发乱序写
			streamTmp := make([]TestMeta, defaultStreamLen)
			copy(streamTmp, TestDataStream)
			streamCnt := defaultStreamLen
			for ; streamCnt > 0; streamCnt-- {
				randId := rand.Intn(streamCnt)
				buf.WriteData([]DataMeta{{streamTmp[randId].data, "0", streamTmp[randId].id,
					streamTmp[randId].status, DataText, "", "", nil}})
				streamTmp = append(streamTmp[:randId], streamTmp[randId+1:]...)
			}
			wg.Done()
		}()
	}

	// 数据流"1"尾数据写入
	buf.WriteData([]DataMeta{{TestDataStream2[defaultStreamLen-1].data, "1", TestDataStream2[defaultStreamLen-1].id,
		TestDataStream2[defaultStreamLen-1].status, DataText, "", "", nil}})
	wg.Wait()

	// 数据流"0"即时读取
	rdTime = time.Now()
	outPut, _, err = buf.ReadMergeData()
	rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if rdPerf > defaultTimeoutDeviation {
		t.Errorf("多数据流并发乱序写入合并读,读耗时过长: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("多数据流并发乱序写入合并读,读操作异常: %s", err.Error())
	}
	// 输出仅len(outPut) == 1, 数据流"0"完整性校验
	if len(outPut) != 1 || bytes.Compare(TestData, outPut[0].Data.([]byte)) != 0 {
		t.Errorf("多数据流场景输入, 输出预期不符, outPut:%v", outPut)
	}

	// 数据流"1"超时读取
	rdTime = time.Now()
	outPut, _, err = buf.ReadMergeData()
	rdPerf = time.Now().Sub(rdTime).Nanoseconds() / (1000 * 1000)
	if math.Abs(float64(rdPerf-defaultTimeout)) > defaultTimeoutDeviation {
		t.Errorf("多数据流并发写入,数据流'1'为超时降级: %d ms", rdPerf)
	} else if err != nil {
		t.Errorf("多数据流并发写入,数据流'1'读操作异常: %s", err.Error())
	}
	// 输出数据校验
	if len(outPut) != 1 || outPut[0].Status != DataStatusLast || outPut[0].FrameId != defaultStreamLen-1 ||
		bytes.Compare(outPut[0].Data.([]byte), TestDataStream2[defaultStreamLen-1].data) != 0 {
		t.Errorf("多数据流并发写入,数据流'1'输出数据异常, outPut: %v", outPut)
	}

	buf.Fini()
	return
}
