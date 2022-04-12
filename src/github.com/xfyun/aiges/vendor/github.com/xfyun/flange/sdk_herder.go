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
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	SDK_METRIC_KAFKA_TOPIC_KEY = "sdk-herder"
)

var (
	// meta data

	// configuration

	// runtime metric, global & consumer level; [outer-code]
	// flush to consumer speed
	flushSuccessGauge int64 = 0
	// flush failed but spill speed
	flushSpillGauge int64 = 0
 	// flush failed to drop speed
	flushFailedGauge int64 = 0
	// consumers's consumerSendGauge
	consumerSendGauge int64 = 0

	// error msg; [outer-code]
	runtimeErrorInterpolation string = ""
)

var (
	// flume client to send sdk-herder to remote trocar
	herderFlumeClient *flumeClient
	// tmp flush host & port
	host = ""
	port = ""

	// header reporter interval
	HerderInterval = 1
)

func initSDKHeader(flumeHost string, flumePort string) {
	host = flumeHost
	port = flumePort

	herderFlumeClient = &flumeClient{
		host: flumeHost,
		port: flumePort,
	}

	go repoterHerder()
}

func repoterHerder() {
	defer catch("reporterHerder")

	herderFlumeClient.open()

	for {
		// reset speed gauge and error str
		atomic.StoreInt64(&flushSuccessGauge, 0)
		atomic.StoreInt64(&flushSpillGauge, 0)
		atomic.StoreInt64(&flushFailedGauge, 0)
		atomic.StoreInt64(&consumerSendGauge, 0)
		runtimeErrorInterpolation = ""

		// sleep for value compute
		time.Sleep(time.Duration(HerderInterval) * time.Second)

		// 1. meta data
		var metaStr = fmt.Sprintf(`"meta": {"schema": "%s", "topic": "%s", "version": "%s"}`,
			DEFAULT_SCHEMA_NAME_VALUE, MESSAGE_QUEUE_KAFKA_TOPIC_VALUE, FLANGE_VERSION)

		// 2. parameter
		var trace = fmt.Sprintf(`"trace": { "deliver_enable": %v, "force_deliver": %v, "dump_enable": %v, "spill_enable": %v, "flush_retry_count": %d, "consumer_num": %d }`,
			DeliverEnable, ForceDeliver, DumpEnable, SpillEnable, FlushRetryCount, consumerClientsNum)
		var consumer = fmt.Sprintf(`"consumer": { "buffer_size": %d, "batch_size": %d, "linger": %d, "flume_host": "%s", "flume_port": "%s" }`,
			BuffSize, BatchSize, LingerSec, host, port)
		var span = fmt.Sprintf(`"svc": { "bcluster": "%s", "idc": "%s", "ip": "%s", "port": "%s", "svn": "%s"}`,
			onceBCluster, onceIDC, onceIPRune, oncePort, onceServiceName)
		var paramterStr = fmt.Sprintf(`"parameter": { %s, %s, %s }`,
			trace, consumer, span)

		// 3. metric
		var metricStr = fmt.Sprintf(`"flush": { "success_gauge": %d, "spill_gauge": %d, "failed_gauge": %d }`,
			atomic.LoadInt64(&flushSuccessGauge), atomic.LoadInt64(&flushSpillGauge), atomic.LoadInt64(&flushFailedGauge))

		// 4. consumer
		channelCap := 0
		channelLen := 0
		for _, c := range consumerClients {
			channelCap += cap(c.chSpanBuffer)
			channelLen += len(c.chSpanBuffer)
		}
		var consumerStr = fmt.Sprintf(`"consumers": { "channel_cap": %d, "channel_len": %d, "send_gauge": %d }`, channelCap, channelLen, atomic.LoadInt64(&consumerSendGauge))

		// 5. error
		var errorStr = fmt.Sprintf(`"error": "%s"`, runtimeErrorInterpolation)

		var mix = "{" + metaStr + "," + paramterStr + "," + consumerStr + "," + metricStr + "," + errorStr + "}"

		// debugf("mix=%s", mix)

		// send mix str to flume
		ev := &flume.ThriftFlumeEvent{Headers: make(map[string]string, 8)}
		ev.Headers[SCHEMA_NAME_KEY] = SDK_METRIC_KAFKA_TOPIC_KEY
		ev.Headers[MESSAGE_QUEUE_KAFKA_TOPIC_KEY] = SDK_METRIC_KAFKA_TOPIC_KEY
		ev.Headers[MESSAGE_QUEUE_RECORD_KEY_KEY] = strconv.Itoa(int(CurrentTimeMillis()))
		ev.Body = str2bytes(mix)

		// make sure rpc lient is open before send msg
		if !herderFlumeClient.rpcClient.Transport.IsOpen() {
			errorf("herder flume client not open, now reset it")
			herderFlumeClient.close()
			herderFlumeClient.transport = nil
			herderFlumeClient.rpcClient = nil
			herderFlumeClient.open()
			continue
		}

		// send
		rs, err := herderFlumeClient.rpcClient.Append(ev)
		if err != nil {
			errorf("herder append error : %v\n", err)
			herderFlumeClient.close()
			herderFlumeClient.transport = nil
			herderFlumeClient.rpcClient = nil
			herderFlumeClient.open()
		} else {
			if rs != flume.Status_OK {
				debugf("herder send 1 items to flume failed.")
			} else {
				debugf("herder send 1 items to flume successfully.")
			}
		}
	}
}
