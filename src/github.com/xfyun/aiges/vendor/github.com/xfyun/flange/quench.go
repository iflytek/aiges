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
	"github.com/golang/protobuf/proto"
	"strconv"
)

var (
	// global endpoint buffer for avro serialize
	endpointString string
)

// serialize spans into one buf
func SerializeSpans(spans []*Span) ([]byte, error) {
	defer catch("serialize error")

	spanBatch := &SpanBatch{
		Endpoint: endpointString,
	}

	// TODO take slice index as span inst, notice empty data,
	// which don't update slice index, cause out-of-order/slice-index-mapping
	// NOTICE: annos & tags start with index [1:] for non-empty confirmed
	for _, span := range spans {
		spanBatch.TraceIds = append(spanBatch.TraceIds, bytes2str(span.traceId))
		spanBatch.Names = append(spanBatch.Names, span.name)
		spanBatch.Ids = append(spanBatch.Ids, bytes2str(span.spanIdTs)+bytes2str(span.spanIdHierarchy))
		spanBatch.Timestamps = append(spanBatch.Timestamps, strconv.Itoa(int(span.timestamp)))
		spanBatch.Durations = append(spanBatch.Durations, strconv.Itoa(int(span.duration)))
		// start [1:0] to eat unused chars
		spanBatch.Annotations = append(spanBatch.Annotations, span.annos[1:])
		// start [1:0] to eat unused chars
		spanBatch.Tags = append(spanBatch.Tags, span.tags[1:])
	}

	return proto.Marshal(spanBatch)
}

// RetrieveSpanInfo with specificed serialize
func RetrieveSpanInfo(data []byte) (traceId string, spanId string, spanType string, sBuf []byte) {
	// NOTICE data with format [c|s|0][len(spanId)spanId][traceId][pb_bytes]
	// get span id
	idLen := uint(data[1])
	spanId = string(data[2 : idLen+2])
	// get trace id
	traceId = string(data[idLen+2 : idLen+34])
	// get span serialize byte
	buf := data[idLen+34:]

	return traceId, spanId, string(data[0]), buf
}