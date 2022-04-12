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
	"github.com/xfyun/flume"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// consumer cache channel size
	BuffSize int32 = 100000
	// log send batch size.
	// The plunger will attempt to batch records together into fewer requests.
	// This helps performance on both the client and the trocar server.
	// A small batch size will make batching less common and may reduce throughput
	// A very large batch size may use memory a bit more wastefully and get over framesize error.
	BatchSize int = 1000
	// linger micro seconds.
	// if have fewer than many records accumulated for batch sends, plunger will 'linger'
	// for the specified time waiting for more records to show up.
	LingerSec int = 5
	// kafka topic for tracar kafka sink routing
	Topic string = MESSAGE_QUEUE_KAFKA_TOPIC_VALUE
)

type consumer struct {
	// flume client to send record to remote trocar
	flumeClient *flumeClient
	// span cache channel
	chSpanBuffer chan *Span
	// span batch for sending
	spanBatches [] *Span

	// exit channel to exit deal-loop sending
	exitChan chan int
	// wait for buffer drain
	wg *sync.WaitGroup
	// wait timer
	waitTicker *time.Ticker

	// flume client index for identify consumer client
	consumerIndex int
	// last deliver timestamp for compute linger time
	lastDeliveryTs int64
}

// init consumer
func initConsumer(flumeHost string, flumePort string, index int, wg *sync.WaitGroup) *consumer {
	// create a new consumer
	c := &consumer{}
	// init flume client
	if flumeHost == "" || flumePort == "" {
		infof("no flume agent address provided, will use default value - %d.", c.consumerIndex)
		flumeHost = FLUME_DEFAULT_HOST
		flumePort = FLUME_DEFAULT_PORT
	}
	c.flumeClient = &flumeClient{
		host: flumeHost,
		port: flumePort,
	}
	// init buffer
	c.chSpanBuffer = make(chan *Span, BuffSize)
	// init span batch
	c.spanBatches = make([] *Span, 0, BatchSize)

	// init exit channel and wg
	c.exitChan = make(chan int, 1)
	c.wg = wg
	c.waitTicker = time.NewTicker(time.Second * time.Duration(LingerSec))

	// assign consumer index and delivery ts
	c.consumerIndex = index
	c.lastDeliveryTs = CurrentTimeMillis()

	// start backend send processing go routine
	c.wg.Add(1)
	go c.appendDataGoroutine()

	return c
}

// stop sending go routine
func (c *consumer) stop() {
	c.waitTicker.Stop()
	c.exitChan <- 1
	infof("stop consumer append go routine - %d", c.consumerIndex)
}

func (c *consumer) appendDataGoroutine() {
	defer catch("append data error")
	defer c.wg.Done()

	infof("start consumer append go routine - %d", c.consumerIndex)

	// start thrift connect
	c.flumeClient.open()
	c.lastDeliveryTs = CurrentTimeMillis()

	// deal loop for sending, exit on channel sign
	for {
		select {
		// sleep while no data
		case <-c.waitTicker.C:
			if DeliverEnable {
				if len(c.spanBatches) == 0 {
					infof("sleep %d second, waiting for span arriving - %d.", LingerSec, c.consumerIndex)
					continue
				} else {
					infof("sleep %d second, with num = %d - %d.", LingerSec, len(c.spanBatches), c.consumerIndex)
					c.sendData()
				}
			}
			// append data to batches
		case tSpan := <-c.chSpanBuffer:
			c.spanBatches = append(c.spanBatches, tSpan)
			// check exit channel sign
		case <-c.exitChan:
			infof("exit channel receive sign - %d.", c.consumerIndex)
			goto FiniDo
		}

		// check connect and batch to send
		if len(c.spanBatches) >= BatchSize || CurrentTimeMillis()-c.lastDeliveryTs > int64(LingerSec*1000) {
			c.sendData()
		}
	}

FiniDo:
	// exit deal-loop send processing
	select {
	case tSpan := <-c.chSpanBuffer:
		c.spanBatches = append(c.spanBatches, tSpan)

		if len(c.spanBatches) >= BatchSize || CurrentTimeMillis()-c.lastDeliveryTs > int64(LingerSec*1000) {
			c.sendData()
		}
	default:
		// nothing to do but exit
	}

	for len(c.spanBatches) > 0 {
		c.sendData()
	}

	if c.flumeClient.transport != nil {
		c.flumeClient.transport.Close()
	}

	infof("exit consumer append go routine : %d", c.consumerIndex)
}

func (c *consumer) sendData() {
	defer catch("sendData")
	begin := CurrentTimeMillis()

	bound := len(c.spanBatches)
	if bound <= 0 {
		debugf("current bound is 0, will jump over this send.")
		return
	}

	// select spans.0 to compute the kafka record key
	kSpan := c.spanBatches[0]
	var rKey []byte
	rKey = append(rKey, kSpan.traceId...)
	rKey = append(rKey, '#')
	rKey = append(rKey, kSpan.spanIdTs...)
	rKey = append(rKey, kSpan.spanIdHierarchy...)
	rKey = append(rKey, '#')
	switch kSpan.spanType {
	case CLIENT, PRODUCER:
		rKey = append(rKey, 'c')
	case SERVER, CONSUMER:
		rKey = append(rKey, 's')
	default:
		rKey = append(rKey, '0')
	}
	rKey = append(rKey, '#')
	rKey = strconv.AppendInt(rKey, int64(bound), 10)

	// serialize spans batch into one buf
	buf, err := SerializeSpans(c.spanBatches)
	// reset span batch for oom
	c.spanBatches = c.spanBatches[bound:]
	if err != nil {
		errorf("send log serialize error : %v.", err)
		return
	}

	// prepare flume event data
	ev := &flume.ThriftFlumeEvent{Headers: make(map[string]string, 8)}
	ev.Headers[MESSAGE_QUEUE_KAFKA_TOPIC_KEY] = Topic
	ev.Headers[SCHEMA_VERSION_KEY] = DEFAULT_SCHEMA_VERSION_VALUE
	ev.Headers[SCHEMA_NAME_KEY] = DEFAULT_SCHEMA_NAME_VALUE
	ev.Headers[SPAN_SERIALIZATION] = "false"
	ev.Headers[TIMESTAMP] = string(kSpan.traceId[8:21])
	ev.Headers[MESSAGE_QUEUE_RECORD_KEY_KEY] = string(rKey)
	ev.Body = buf

	// make sure rpc client is open before send msg
	if !c.flumeClient.rpcClient.Transport.IsOpen() {
		errorf("flume client not open, now reset it")
		c.flumeClient.close()
		c.flumeClient.transport = nil
		c.flumeClient.rpcClient = nil
		c.flumeClient.open()
		return
	}

	// send
	rs, err := c.flumeClient.rpcClient.Append(ev)
	if err != nil {
		errorf("append batch error : %v\n", err)
		c.flumeClient.close()
		c.flumeClient.transport = nil
		c.flumeClient.rpcClient = nil
		c.flumeClient.open()
	} else {
		if rs != flume.Status_OK {
			debugf("cid:%d send %d spans to flume failed, will try again later.", c.consumerIndex, bound)
		} else {
			atomic.AddInt64(&consumerSendGauge, int64(bound))
			c.lastDeliveryTs = CurrentTimeMillis()
			debugf("cid:%d send %d items to flume successfully, use %d ms.", c.consumerIndex, bound, c.lastDeliveryTs-begin)
		}
	}
}
