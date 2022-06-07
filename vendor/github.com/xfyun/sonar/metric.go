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
	"fmt"
	"os"
	"errors"
)

// create a new metric
func newMetric(metricType string, metricName string, endpoint string) *MetricData {
	return &MetricData{
		endpoint:endpoint,
		metric:metricName,
		counterType:metricType,
		step:10,
		timestamp:CurrentTimeMillis() / 1000,
		tagsMap:make(map[string]string)}
}

// create a new metric with serverName and port
func NewMetricWithNamePort(metricType string, metricName string, endpoint string, serviceName string, port string, ds string) *MetricData {
	var metric = newMetric(metricType, metricName, endpoint)
	metric.tagsMap["svc"] = serviceName
	metric.tagsMap["port"] = port
	metric.tagsMap["ds"] = ds
	return metric
}

// create a new metric with global value set
func NewMetricWithGlobal(metricType string, metricName string) *MetricData {
	// validate the global value
	if ( len(Endpoint) > 0 && len(ServiceName) > 0 && len(Port) > 0 && len(DS) > 0) {
		var metric = newMetric(metricType, metricName, Endpoint)
		metric.tagsMap["svc"] = ServiceName
		metric.tagsMap["port"] = Port
		metric.tagsMap["ds"] = DS
		return metric
	}
	// if validate failed
	return nil
}

// set step
func (metric *MetricData) WithStep(step int64) *MetricData {
	metric.step = step
	return metric
}

// set value
func (metric *MetricData) WithValue(value float64) *MetricData {
	metric.value = value
	return metric
}

// set tag
func (metric *MetricData) Tag(tag KV) *MetricData {
	var value = ""
	switch tag.Value.(type) {
	case bool:
		if tag.Value == true { value = "true" } else { value = "false"}
	case int, int64, int32, int16, int8, uint8, uint16, uint32, uint64, uint:
		value = fmt.Sprintf("%d", tag.Value)
	case float32, float64:
		value = fmt.Sprintf("%f", tag.Value)
	case string:
		value = fmt.Sprintf("%s", tag.Value)
	default:
		panic("unsupport type")
	}
	metric.tagsMap[tag.Key] = value
	return metric
}

// tag ds
func (metric *MetricData) TagDS(ds string) *MetricData {
	return metric.Tag(KV{"ds", ds})
}

// convert to string in json
func (metric *MetricData) ToString() string {
	tags := ""
	for k,v := range metric.tagsMap {
		tags += k + "=" + v + ","
	}
	return fmt.Sprintf("{\"endpoint\":\"%s\", \"metric\":\"%s\", \"timestamp\":%d, " +
		"\"value\":%f, \"counterType\":\"%s\", \"step\":\"%d\", \"tags\":\"%s\"}", metric.endpoint, metric.metric, metric.timestamp, metric.value, metric.counterType, metric.step, tags[0:len(tags)-1])
}

// dump span to file
// ${dir}/${endpoint}_${metric}
func (metric *MetricData) dump(dir string) {
	filename := fmt.Sprintf("%s%s%s_%s_%d", dir, string(os.PathSeparator), metric.endpoint, metric.metric, metric.timestamp)
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		if Logger != nil { Logger.Debugf("open %s failed: %v", filename, err) }
		return
	}

	// Dump span
	fmt.Fprintln(fp, metric.ToString())
	fp.Close()
}

// spill data while channel is full
// it will spill data to local disk for further process
func (metric *MetricData) spill() {
	// fixed directory, which means uniform process
	spillDir := "." + string(os.PathSeparator) + "log" + string(os.PathSeparator) + "metric_spill"
	if err := os.MkdirAll(spillDir, 0755); err != nil {
		if Logger != nil { Logger.Debugf("mkdir dumpdir err: %v", err) }
		return
	}
	// spill just like dump
	metric.dump(spillDir)
}

// flush dslog to buffer channel, will be send to flume.
func (metric *MetricData) Flush() error {
	defer catch("flush")

	// dump to file
	if DumpEnable {
		metric.dump(DumpDir)
	}

	// save to buffer only while plunger is enabled
	if DeliverEnable {
		select {
		case chSpanBuff <- metric:
			//log.Println("Produce one log to buffer.")
			return nil
		default:
			// if timeout, then spill this value to disk for further process
			metric.spill()
			if Logger != nil { Logger.Debugf("Save one log to buffer timeout")}
			return errors.New("Save one log to buffer timeout")
		}
	}

	return nil
}

