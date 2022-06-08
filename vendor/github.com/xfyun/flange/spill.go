/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package flange

import (
	"errors"
	"github.com/xfyun/flume"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// log spill dir
	SpillDir string = "." + string(os.PathSeparator) + "spill"
	// flag to enable spill
	SpillEnable bool = true
	// max spill content size in G-byte
	MaxSpillContentSize int64 = 1
)

var (
	// flume client
	reverseSpillFlumeClient *flumeClient
	// temp flumeEvent batch
	reverseSpillEvent *flume.ThriftFlumeEvent

	// spillRingBuffer for spill data in-memory saving
	chSpillBuffer chan *Span

	// spill date-time name, segmented by `yyyy-mm-dd-hh`
	spillSegmentDateName string
	// spill data & index file name
	spillDataFile *os.File = nil
	spillIndexFile *os.File = nil
	// current spill content size
	currentContentSize int64

	// chan to control exit go routine
	spillExitChan chan int
	// wait for buffer drain
	spillWg *sync.WaitGroup
)

// init spill profile
func initSpillProf(flumeHost string, flumePort string, wg *sync.WaitGroup) error {
	// check & create spill dir
	if SpillDir != "" {
		debugf("mkdir spill dir")
		if err := os.MkdirAll(SpillDir, 0755); err != nil {
			errorf("mkdir spillDir err : %v", err)
			return err
		}
	} else {
		errorf("unvalid spill dir : %v", SpillDir)
		return errors.New("unvalid spill dir : " + SpillDir)
	}

	// exit wait group
	spillExitChan = make(chan int)
	spillWg = wg

	// cache channel buffer
	chSpillBuffer = make(chan *Span, BuffSize)

	spillSegmentDateName = time.Now().Format("2006-01-02-15")
	var err error
	// init spill data file
	spillDataName := SpillDir + string(os.PathSeparator) +  "span.spill." + spillSegmentDateName + ".data"
	spillDataFile, err = os.OpenFile(spillDataName, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		errorf("open spill data file %s failed : %v", spillDataName, err)
		return err
	} else {
		infof("create spill data file success, now seek it to end")
		spillDataFile.Seek(0, os.SEEK_END)
	}
	// init spill index file
	spillIndexName := SpillDir + string(os.PathSeparator) + "span.spill." + spillSegmentDateName + ".idx"
	spillIndexFile, err = os.OpenFile(spillIndexName, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		errorf("open spill index file %s failed : %v", spillIndexName, err)
		return err
	} else {
		infof("create spill index file success, now seek it to send")
		spillIndexFile.Seek(0, os.SEEK_END)
	}

	// init flume client
	reverseSpillFlumeClient = &flumeClient{
		host: flumeHost,
		port: flumePort,
	}
	reverseSpillEvent = &flume.ThriftFlumeEvent{Headers: make(map[string]string, 8)}
	reverseSpillEvent.Headers[MESSAGE_QUEUE_KAFKA_TOPIC_KEY] = Topic
	reverseSpillEvent.Headers[SCHEMA_VERSION_KEY] = DEFAULT_SCHEMA_VERSION_VALUE
	reverseSpillEvent.Headers[SCHEMA_NAME_KEY] = DEFAULT_SCHEMA_NAME_VALUE
	reverseSpillEvent.Headers[SPAN_SERIALIZATION] = "false"

	spillWg.Add(1)
	go spillReverseGoroutine()

	return nil
}

// stop current spill go routine
func stopSpill() {
	if spillExitChan != nil {
		spillExitChan <- 1
	}
}

// start spill and reverse send go routine
func spillReverseGoroutine() {
	defer catch("spill reverse go routine error")
	defer spillWg.Done()

	infof("start spill & reverse go routine")

	// seek spill fd
	currentContentSize, _ = spillDataFile.Seek(0, io.SeekEnd)
	// init reverse spill flume client
	reverseSpillFlumeClient.open()

	for {
		// retrieve a span\timeout\exit
		select {
		// check spill channel
		case tSpan := <-chSpillBuffer:
			if currentContentSize < MaxSpillContentSize*1024*1024*1024 {
				spillSpan(tSpan)
			}
			// no spill span, check reverse file
		case <-time.After(time.Second * time.Duration(LingerSec)):
			files, err := ioutil.ReadDir(SpillDir + string(os.PathSeparator))
			if err != nil {
				errorf("read spill dir files error with %v", err)
			} else {
				for _, f := range files {
					if strings.Contains(f.Name(), ".data") && !strings.Contains(f.Name(), spillSegmentDateName) {
						readBack(f.Name())
					}
				}
			}
			// exit signal
		case <-spillExitChan:
			goto FiniSpill
		}
	}

FiniSpill:
	debugf("exit spill go routine")
	select {
	case tSpan := <-chSpillBuffer:
		if currentContentSize < MaxSpillContentSize*1024*1024*1024 {
			spillSpan(tSpan)
		}
	default:
		// nothing to do
	}

	// exit
	if reverseSpillFlumeClient != nil {
		reverseSpillFlumeClient.close()
	}
	if spillDataFile != nil {
		spillDataFile.Close()
	}
	if spillIndexFile != nil {
		spillIndexFile.Close()
	}
	infof("exit spill & reverse go routine")
}

// spill span to local file
func spillSpan(span *Span) {
	defer catch("spillSpan")

	buf, err := SerializeSpans([]*Span{span})
	if err != nil {
		errorf("spill serialize error : %v", err)
		return
	}
	// add span_meta_info for reverse spill
	var spillBuf []byte
	// span type
	switch span.spanType {
	case CLIENT, PRODUCER:
		spillBuf = append(spillBuf, 'c')
	case SERVER, CONSUMER:
		spillBuf = append(spillBuf, 's')
	default:
		spillBuf = append(spillBuf, '0')
	}
	// idl
	idLen := len(span.spanIdTs) + len(span.spanIdHierarchy)
	spillBuf = append(spillBuf, byte(idLen))
	spillBuf = append(spillBuf, span.spanIdTs...)
	spillBuf = append(spillBuf, span.spanIdHierarchy...)
	// traceId
	spillBuf = append(spillBuf, span.traceId...)
	// all
	spillBuf = append(spillBuf, buf...)

	// check if `spillSegmentDateName` changed
	spillSegmentDateName = time.Now().Format("2006-01-02-15")
	if !strings.Contains(spillDataFile.Name(), spillSegmentDateName) {
		// rebuild `spillDataFile`
		spillDataFile.Close()
		spillDataName := SpillDir + string(os.PathSeparator) +  "span.spill." + spillSegmentDateName + ".data"
		spillDataFile, err = os.OpenFile(spillDataName, os.O_RDWR | os.O_CREATE, 0666)
		if err != nil {
			errorf("open spill data file %s failed : %v", spillDataName, err)
		} else {
			infof("update spill data file - [%s] success, now reset curContSize and seek it to end", spillDataName)
			spillDataFile.Seek(0, os.SEEK_END)
			currentContentSize = 0
		}

		// rebuild `spillIndexFile`
		spillIndexFile.Close()
		spillIndexName := SpillDir + string(os.PathSeparator) + "span.spill." + spillSegmentDateName + ".idx"
		spillIndexFile, err = os.OpenFile(spillIndexName, os.O_RDWR | os.O_CREATE, 0666)
		if err != nil {
			errorf("open spill index file %s failed : %v", spillIndexName, err)
		} else {
			infof("update spill index file - [%s] success, now seek it to end", spillIndexName)
			spillIndexFile.Seek(0, os.SEEK_END)
		}
	}

	if strings.Contains(spillDataFile.Name(), spillSegmentDateName) && strings.Contains(spillIndexFile.Name(), spillSegmentDateName) {
		// spill data with index
		l, err := fmt.Fprint(spillDataFile, string(spillBuf))
		if err != nil {
			errorf("spill data error with %v", err)
			return
		}
		currentContentSize += int64(l)
		_, err = fmt.Fprintf(spillIndexFile, "%d\n", l)
		if err != nil {
			errorf("spill index error with %v", err)
			// if index file write failed, then roll back for data file
			spillDataFile.Seek(-1 * int64(l), os.SEEK_END)
		}
	} else {
		errorf("spill with diff data name - [%s] and index name - [%s]", spillDataFile.Name(), spillIndexFile.Name())
	}
}

func readBack(dataName string) {
	defer catch("readBack")

	debugf("current read back file name = " + dataName)

	spillDataName := SpillDir + string(os.PathSeparator) + dataName
	defer os.Remove(spillDataName)
	spillIndexName := SpillDir + string(os.PathSeparator) + strings.Replace(dataName, ".data", ".idx", -1)
	defer os.Remove(spillIndexName)

	idxBuf, e := ioutil.ReadFile(spillIndexName)
	if e != nil {
		errorf("read index file failed %v", e)
		return
	}

	spillDataFile, err := os.OpenFile(spillDataName, os.O_RDWR | os.O_CREATE, 0666)
	defer spillDataFile.Close()
	if err != nil {
		errorf("read data file failed %v", err)
		return
	}

	for _, n := range strings.Split(string(idxBuf), "\n") {
		if len(n) <= 0 {
			continue
		}

		nInt, e := strconv.Atoi(n)
		if e != nil {
			errorf("read index value failed %v", err)
			return
		}

		lengthBuffer := make([]byte, nInt)
		if l, err := spillDataFile.Read(lengthBuffer); l == nInt && err == nil {
			sendLog(lengthBuffer)
		} else {
			errorf("read data file failed %v", err)
			return
		}
	}
}

func sendLog(bytes []byte) bool {
	defer catch("reverse spill send log")

	// retrieve span info
	traceId, spanId, spanType, buf := RetrieveSpanInfo(bytes)
	infof("retrieve tid=%s, sid=%s", traceId, spanId)

	reverseSpillEvent.Body = buf
	// set headers, see java syringe.v2 FlumeClient
	// set timestamp to the millseconds of the start of this trace
	reverseSpillEvent.Headers[TIMESTAMP] = traceId[8:21]
	// set r.k sid for mq balancing, format: <traceId>
	reverseSpillEvent.Headers[MESSAGE_QUEUE_RECORD_KEY_KEY] = traceId + "#" + spanId + "#" + spanType + "#1"
	// fc.evBatches[i].Headers[FLUSH_TIMESTAMP] = strconv.FormatInt(CurrentTimeMillis(), 10)

	// make sure rpc client is open before send msg
	if !reverseSpillFlumeClient.rpcClient.Transport.IsOpen() {
		errorf("reverse spill flume client is not open, now reset it")
		reverseSpillFlumeClient.close()
		reverseSpillFlumeClient.transport = nil
		reverseSpillFlumeClient.rpcClient = nil
		reverseSpillFlumeClient.open()
		return false
	}

	rs, err := reverseSpillFlumeClient.rpcClient.Append(reverseSpillEvent)
	if err != nil {
		errorf("reverseSpill send 1 log to flume failed, will try again later.")
		reverseSpillFlumeClient.close()
		reverseSpillFlumeClient.transport = nil
		reverseSpillFlumeClient.rpcClient = nil
		reverseSpillFlumeClient.open()
		return false
	} else {
		if rs != flume.Status_OK {
			debugf("reverseSpill send 1 log to flume failed, will try again later.")
			return false
		} else {
			debugf("reverseSpill send 1 log successful.")
			return true
		}
	}
}
