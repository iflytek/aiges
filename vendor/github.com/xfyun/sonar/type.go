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

// CounterType
const (
	// metric counter type for `gauge`
	TYPE_GAUGE = "GAUGE"
	// metric counter type for `counter`
	TYPE_COUNTER = "COUNTER"
)

// RPC settings
const (
	// default flume agent host
	FLUME_DEFAULT_HOST = "127.0.0.1"
	// default flume agent port
	FLUME_DEFAULT_PORT = "4545"
)

// Flume event headers' constants
const (
	// timestamp, need set to flume event, CAN NOT USE `ts`
	TIMESTAMP = "timestamp"
	// session id
	SID = "sid"
	// schema version key
	SCHEMA_VERSION_KEY = "s.v"
	// schema name key
	SCHEMA_NAME_KEY = "s.n"
	// division name key
	DIVISION_NAME_KEY = "d.n"
	// project name key
	PROJECT_NAME_KEY = "p.n"
	// default schema version
	DEFAULT_SCHEMA_VERSION_VALUE = "1"
	// default schema name
	DEFAULT_SCHEMA_NAME_VALUE = "vagus"
	// key of message queue(kafka) topic
	MESSAGE_QUEUE_KAFKA_TOPIC_KEY = "k.t"
	// message queue(kafka) default topic
	MESSAGE_QUEUE_KAFKA_TOPIC_VALUE = "vagus"
	// key of message queue key
	MESSAGE_QUEUE_RECORD_KEY_KEY = "r.k"
	// flush timestamp
	FLUSH_TIMESTAMP = "flush.ts"
)

// metric data to collect
type MetricData struct {
	// server/component deploy endpoint
	endpoint string
	// current metric name to collect
	metric string
	// metric create timestamp
	timestamp int64
	// the metric value
	value float64
	// metric type , such as [gauge, counter]
	counterType string
	// metric collect step, every `step` to do
	step int64
	// tags, extra message to save
	tags string
	// tags map
	tagsMap map[string]string
}

// metric tags k-v pair
type KV struct {
	// tag key
	Key string
	// tag value
	Value interface{}
}

// custom log interface
type CustomLogInterface interface {

	// Infof formats message according to format specifier
	// and writes to log with level = Info.
	Infof(format string, params ...interface{})

	// Debugf formats message according to format specifier
	// and writes to log with level = Debug.
	Debugf(format string, params ...interface{})
}