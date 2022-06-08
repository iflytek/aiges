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

package sonar

import (
	"os"
	"time"
	"sync"
	"sync/atomic"
	"errors"
)

var (
	// flume clients
	clients []FlumeClient = nil
	// default flume client number
	FlumeClientNum = 4
	// span buffer
	chSpanBuff chan *MetricData = nil
	// wait for buffer drain
	wg *sync.WaitGroup
	// global init check with cas
	atomicInitInteger int32 = 0
)

var (
	// log buffer size.
	BuffSize int32 = 20000
	// log send batch size.
	// The plunger will attempt to batch records together into fewer requests.
	// This helps performance on both the client and the trocar server.
	// A small batch size will make batching less common and may reduce throughput
	// A very large batch size may use memory a bit more wastefully and get over framesize error.
	BatchSize int = 100
	// linger seconds.
	// if have fewer than many records accumulated for batch sends, plunger will 'linger'
	// for the specified time waiting for more records to show up.
	LingerSec int = 5
	// kafka topic for tracar kafka sink routing
	Topic string = MESSAGE_QUEUE_KAFKA_TOPIC_VALUE
	// true if dump log to file.
	DumpEnable bool = false
	// log dump dir.
	DumpDir string = "." + string(os.PathSeparator) + "metric"
	// true if deliver log to flume.
	DeliverEnable bool = true
	// logger
	Logger CustomLogInterface

	// global endpoint value
	Endpoint string = ""
	// global service name value
	ServiceName string	= ""
	// global port value
	Port string = ""
	// global data source value
	DS string = ""
)

// initialize plunger.
// @param host:port determine the flume host:ip
// @param num number of backend-data-consumer in Flange, adjust with upstream load
func Init(host string, port string, num int) error {
	// check init already
	if !atomic.CompareAndSwapInt32(&atomicInitInteger, 0, 1) {
		return errors.New("already init, ignore this")
	}

	consumerNum := num
	if num > 10 || num < 1 {
		consumerNum = FlumeClientNum
	}

	// mkdir dump dir
	if DumpEnable && DumpDir != "" {
		if Logger != nil { Logger.Infof("mkdir dumpdir") }
		if err := os.MkdirAll(DumpDir, 0755); err != nil {
			if Logger != nil { Logger.Debugf("mkdir dumpdir err: %v", err) }
			return err
		}
	}

	// init buffer
	chSpanBuff = make(chan *MetricData, BuffSize)
	wg = &sync.WaitGroup{}

	if clients == nil {
		for i := 0; i < consumerNum; i++ {
			client := FlumeClient{
				Host : host,
				Port : port}
			client.init(chSpanBuff, wg)
			clients = append(clients, client)
		}
	}

	return nil
}

// set the global service properties
func SetGlobalProp(endpoint string, serviceName string, port string, ds string) {
	Endpoint = endpoint
	ServiceName = serviceName
	Port = port
	DS = ds
}

// fini plunger.
func Fini() {
	defer catch("finish")

	for {
		select {
		case <- time.After(time.Second * time.Duration(LingerSec)):
			goto Out
		case ev := <- chSpanBuff:
			chSpanBuff <- ev
		}
	}
Out:
	// clear the channel data first
	for _, client := range clients {
		client.Exit <- 1
	}
	// wait for process
	wg.Wait()
	// close the transport
	for _, client := range clients {
		client.transport.Close()
	}
}

// global error catch function
func catch(site string) {
	if err := recover(); err != nil {
		if Logger != nil { Logger.Debugf("Error occur [%v] at [%s]", err, site)}
	}
}

