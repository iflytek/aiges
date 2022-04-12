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
	"github.com/xfyun/thrift"
	"github.com/xfyun/flume"
	"sync"
	"net"
	"errors"
	"time"
	"reflect"
	"fmt"
	"strconv"
)

// flume thrift rpc client
type FlumeClient struct {
	// flume agent host
	Host string
	// flume agent port
	Port string
	// exit tag
	Exit chan int

	// span buffer
	chMetricBuffer chan *MetricData
	// true if init successed
	initSucc bool

	// last deliver timestamp
	lastDeliveryTs int64

	// span queue for serialize and append to flume in batch
	metricQueue [] *MetricData
	// socket
	transport *thrift.TSocket
	// flume sdk client
	rpcClient *flume.ThriftSourceProtocolClient

	// wait for buffer drain
	wg *sync.WaitGroup
}

// initialize flume client
func (fc *FlumeClient) init(metricBuffer chan *MetricData, wg *sync.WaitGroup) error {
	fc.initSucc = false

	// init buffer
	fc.chMetricBuffer = metricBuffer

	if fc.Host == "" || fc.Port == "" {
		if Logger != nil { Logger.Infof("No flume agent address provided, will use default") }
		fc.Host = FLUME_DEFAULT_HOST
		fc.Port = FLUME_DEFAULT_PORT
	}

	// prepare
	if err := fc.prepare(); err != nil {
		if Logger != nil { Logger.Debugf("Plunger init err: %v", err) }
		return err
	}

	fc.wg = wg
	fc.wg.Add(1)
	fc.lastDeliveryTs = CurrentTimeMillis()

	fc.Exit = make(chan int, 1)
	go fc.appendMetrics()
	fc.initSucc = true

	return nil
}

// reopen transport
func (fc *FlumeClient) prepare() error {
	// transport settings, see java syringe
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTCompactProtocolFactory()

	var err error
	fc.transport, err = thrift.NewTSocket(net.JoinHostPort(fc.Host, fc.Port))
	if err != nil {
		if Logger != nil { Logger.Debugf("transport err : %v", err) }
		return errors.New("thrfit transport build err")
	}

	fc.transport.SetTimeout(time.Second * 10)
	useTransport := transportFactory.GetTransport(fc.transport)
	if useTransport == nil {
		return errors.New("thrfit not get transport")
	}

	fc.rpcClient = flume.NewThriftSourceProtocolClientFactory(useTransport, protocolFactory)
	if fc.rpcClient == nil {
		return errors.New("thrfit new client failed")
	}

	return nil
}

// fini flume client.
func (fc *FlumeClient) Fini() {
	fc.Exit <- 1
	fc.transport.Close()
}

// true if is flume clinet is opened.
// is transport is closed, will try reopen.
func (fc *FlumeClient) IsOpen() bool {
	if fc.transport.IsOpen() {
		return true
	}

	// reopen
	if err := fc.transport.Open(); err != nil {
		if Logger != nil { Logger.Debugf("%v", err) }
		if reflect.TypeOf(err).String() == "*thrift.tTransportException" {
			if err.(thrift.TTransportException).TypeId() == thrift.ALREADY_OPEN {
				return true
			} else {
				if Logger != nil { Logger.Debugf("%v", err) }
				return false
			}
		}
	}

	return false
}

// append spans from buffer to queue,
// if reach at Batch or time limit then send.
func (fc *FlumeClient) appendMetrics() {
	defer fc.wg.Done()

	if Logger != nil { Logger.Infof("Start append routine") }
	for {
		select {
		case <-time.After(time.Second * time.Duration(LingerSec)):
			if DeliverEnable {
				if len(fc.metricQueue) == 0 {
					if Logger != nil { Logger.Debugf("sleep %d sec, wait for span arriving", LingerSec) }
					continue
				} else {
					if Logger != nil { Logger.Infof("sleep %d sec, with num = %d", LingerSec, len(fc.metricQueue)) }
					fc.sendLogs()
				}
			}
		case evTask := <-fc.chMetricBuffer:
			fc.metricQueue = append(fc.metricQueue, evTask)
		case <-fc.Exit:
			if Logger != nil { Logger.Infof("Exit") }
			goto FiniDo
		}
		if fc.IsOpen() {
			if len(fc.metricQueue) >= BatchSize || CurrentTimeMillis()-fc.lastDeliveryTs > int64(LingerSec*1000) {
				fc.sendLogs()
			}
		} else {
			if Logger != nil { Logger.Debugf("Thrift Server can not connect, retry later") }
			//goto FiniDo
		}
	}
FiniDo:
	if len(fc.metricQueue) != 0 {
		if fc.IsOpen() {
			for len(fc.metricQueue) > 0 {
				fc.sendLogs()
			}
		}
	}
	if Logger != nil { Logger.Infof("Fini") }
}

// append logs from queue to flume agent in batches.
func (fc *FlumeClient) sendLogs() {
	begin := CurrentTimeMillis()
	bound := Min(BatchSize, len(fc.metricQueue))
	batch := fc.metricQueue[:bound]

	var evBatches []*flume.ThriftFlumeEvent
	for _, metric := range batch {
		// build flume event
		ev := &flume.ThriftFlumeEvent{
			Headers:make(map[string]string),
			Body: []byte(metric.ToString())}

		// set headers, see java syringe.v2 FlumeClient
		// set timestamp to the millseconds of the start of this trace
		ev.Headers[TIMESTAMP] =  fmt.Sprintf("%d", metric.timestamp)
		// set k.t for mq(kafka) topic
		ev.Headers[MESSAGE_QUEUE_KAFKA_TOPIC_KEY] = Topic
		// set r.k sid for mq balancing, format: <traceId>
		ev.Headers[MESSAGE_QUEUE_RECORD_KEY_KEY] = Md5(PK(metric.endpoint, metric.metric, metric.tagsMap))
		ev.Headers[SCHEMA_VERSION_KEY] = DEFAULT_SCHEMA_VERSION_VALUE
		ev.Headers[SCHEMA_NAME_KEY] = DEFAULT_SCHEMA_NAME_VALUE
		ev.Headers[FLUSH_TIMESTAMP] = strconv.FormatInt(CurrentTimeMillis(), 10)

		evBatches = append(evBatches, ev)
	}
	rs, err := fc.rpcClient.AppendBatch(evBatches)

	if err != nil {
		// TODO(gwjiang): check overlimit error and turn down the BatchSize
		if Logger != nil { Logger.Debugf("%v", err) }
		// close before reopen
		fc.transport.Close()
		fc.prepare()
		fc.IsOpen()
	} else {
		if rs != flume.Status_OK {
			if Logger != nil { Logger.Infof("send %d spsns to flume failed, will try again later.", len(batch)) }
			if Logger != nil { Logger.Debugf("%v", err) }
		} else {
			// clean batch from logs which has been sent successfully
			if Logger != nil { Logger.Debugf("send %d spans to flume successfully, use %d ms.",
				len(batch), CurrentTimeMillis()-begin) }
			fc.metricQueue = fc.metricQueue[bound:]
			fc.lastDeliveryTs = CurrentTimeMillis()
		}
	}
}