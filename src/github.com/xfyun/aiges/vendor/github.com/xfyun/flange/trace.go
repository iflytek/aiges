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
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var (
	// Deprecated
	// sleep ts in gen under low load by micro-second
	LowLoadSleepTs = 100
	// Deprecated
	// switch for self web based log metric
	WatchLogEnable = true
	WatchLogPort   = 12331
)

var (
	// consumer clients
	consumerClients []*consumer = nil
	// default flume client number
	consumerClientsNum = 8
	// wait for buffer drain
	wg *sync.WaitGroup
)

var (
	// true if dump log to file.
	DumpEnable bool = false
	// log dump dir.
	DumpDir string = "." + string(os.PathSeparator) + "trace"

	// global init check with cas
	atomicInitInteger int32 = 0
	// count drop span
	atomicSpanDropCount int64 = 0
	// small then FlumeClientNum
	balanceIndex int = 0

	// retry count
	FlushRetryCount int = 10
	// true if deliver log to flume.
	DeliverEnable bool = true
	// force deliver the sample trace
	ForceDeliver bool = false
)

// initialize plunger.
// @param host:port determine the flume host:ip
// @param num number of backend-data-consumer in Flange, adjust with upstream load
func Init(flumeHost string, flumePort string, cNum int, bcluster string, idc string, serviceIP string, servicePort string, serviceName string) error {
	// set global value
	setGlobalConfig(bcluster, idc, serviceIP, servicePort, serviceName)

	// check init already
	if !atomic.CompareAndSwapInt32(&atomicInitInteger, 0, 1) {
		infof("already init flange, ignore this operation.")
		return nil
	}

	// mkdir dump dir
	if DumpEnable && DumpDir != "" {
		debugf("mkdir dumpdir.")
		if err := os.MkdirAll(DumpDir, 0755); err != nil {
			errorf("mkdir dumpdir err : %v.", err)
			return err
		}
	}

	wg = &sync.WaitGroup{}
	// init spill info, not start, just wait for mg_center
	initSpillProf(flumeHost, flumePort, wg)
	initSDKHeader(flumeHost, flumePort)

	// uncheck consumer cNum
	if cNum < 1 {
		infof("consumer cNum should not less than 1, using default cNum 8.")
		cNum = consumerClientsNum
	}
	consumerClientsNum = cNum
	for i := 0; i < consumerClientsNum; i++ {
		c := initConsumer(flumeHost, flumePort, i, wg)
		consumerClients = append(consumerClients, c)
	}

	return nil
}

// flush dslog to buffer channel, will be send to flume.
func Flush(span *Span) error {
	defer catch("Flush")

	if span == nil {
		return errors.New("span is nil")
	}

	// check non-empty anno & tags for batch-serialize index confirmed
	if len(span.annos) <= 0 || len(span.tags) <= 0 {
		span.WithTag("invalid", "no anno || no tag")
		return errors.New("invalid span with no anno or no tag")
	}

	// dump to file
	if DumpEnable {
		span.dump(DumpDir)
	}

	// support for sdk sample, 'a' for abandon in trace id at index 25
	if !ForceDeliver && span.traceId[25] == 'a' {
		// release span
		return nil
	}

	// save to buffer only while plunger is enabled
	if DeliverEnable {
		// add span-tree to tree-ring-buffer with concurrent assigned balanceIndex for balance, see `flume_client.go`
		for i := 0; i < FlushRetryCount; i++ {
			// compute balanceIndex each time for fixed value failed
			balanceIndex = int(CurrentTimeMillis() % int64(consumerClientsNum))

			select {
			case consumerClients[balanceIndex].chSpanBuffer <- span:
				atomic.AddInt64(&flushSuccessGauge, 1)
				return nil
			default:
				continue
			}
		}

		// add span to spill
		if SpillEnable {
			select {
			case chSpillBuffer <- span:
				atomic.AddInt64(&flushSpillGauge, 1)
				return nil
			default:
				// spill channel is full
			}
		}

		// print out to drop
		atomic.AddInt64(&flushFailedGauge, 1)
		if atomic.AddInt64(&atomicSpanDropCount, 1) > 10000 {
			infof("%d drop span count = 10000", CurrentTimeMillis())
			atomic.StoreInt64(&atomicSpanDropCount, 0)
		}
	}

	return nil
}

// fini plunger.
func Fini() {
	defer catch("finish")

	// stop all component
	for _, c := range consumerClients {
		c.stop()
	}
	stopSpill()
	wg.Wait()
}

// global error catch function
func catch(site string) {
	if err := recover(); err != nil {
		errorf("Error occur [%v] at [%s] with \n%s.", err, site, string(debug.Stack()))
		// debug.PrintStack()
	}
}
