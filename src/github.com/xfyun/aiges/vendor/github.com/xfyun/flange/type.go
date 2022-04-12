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

// Constants.
const (
	// The client sent ("cs") a request to a server. There is only one send per span.
	CLIENT_SEND = "cs"
	// The client received ("cr") a response from a server. There is only one receive per span.
	CLIENT_RECV = "cr"
	// The server sent ("ss") a response to a client. There is only one response per span.
	SERVER_SEND = "ss"
	// The server received ("sr") a request from a client. There is only one request per span.
	SERVER_RECV = "sr"
	// Message send ("ms") is a request to send a message to a destination, usually a broker.
	MESSAGE_SEND = "ms"
	// A consumer received ("mr") a message from a broker.
	MESSAGE_RECV = "mr"
	// When an Endpoint.Value, this indicates when an error occurred.
	ERROR = "error"
	// The Tag.Value of "lc" is the component or namespace of a local span.
	LOCAL_COMPONENT = "lc"
	// When present, Tag.Endpoint indicates a client address ("ca") in a span.
	CLIENT_ADDR = "ca"
	// When present, Tag.Endpoint indicates a server address ("sa") in a span.
	SERVER_ADDR = "sa"
	// Indicates the remote address of a messaging span, usually the broker.
	MESSAGE_ADDR = "ma"
)

// Type
const (
	// unknown type.
	UNKNOWN int32 = 0
	// for CLIENT_SEND or CLIENT_RECV.
	CLIENT int32 = 1
	// for SERVER_RECV or SERVER_SEND.
	SERVER int32 = 2
	// for MESSAGE_SEND.
	PRODUCER int32 = 3
	// for MESSAGE_RECV.
	CONSUMER int32 = 4
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
	DEFAULT_SCHEMA_NAME_VALUE = "trace"
	// key of message queue(kafka) topic
	MESSAGE_QUEUE_KAFKA_TOPIC_KEY = "k.t"
	// message queue(kafka) default topic
	MESSAGE_QUEUE_KAFKA_TOPIC_VALUE = "trace-v2"
	// key of message queue key
	MESSAGE_QUEUE_RECORD_KEY_KEY = "r.k"
	// flush timestamp
	FLUSH_TIMESTAMP = "flush.ts"
	// serialize options
	SPAN_SERIALIZATION = "span.ser"

	// TODO remember to update this while version update
	// flange version
	FLANGE_VERSION = "1.1.0"
)
