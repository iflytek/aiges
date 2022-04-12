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
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Get millis timestamp.
func CurrentTimeMillis() int64 {
	return int64(time.Now().UnixNano() / int64(time.Millisecond))
}

// Get micros timestamp.
func CurrentTimeMicros() int64 {
	return int64(time.Now().UnixNano() / int64(time.Microsecond))
}

// Get formative datetime string.
func GetDateTimeStrMillis(ts int64) string {
	return time.Unix(0, ts*1000000).Format("2006-01-02 15:04:05.000")
}

// Get formative datetime string.
func GetDateTimeStrMicros(ts int64) string {
	return time.Unix(0, ts*1000).Format("2006-01-02 15:04:05.000")
}

// Convert ip string to int32.
func IPv4toRune(ip string) int32 {
	var ret int32 = 0
	for _, par := range strings.Split(ip, ".") {
		// shift the previously parsed bits over by one byte
		ret = ret << 8
		// set the low order bits to the current octet
		i, _ := strconv.Atoi(par)
		ret |= int32(i)
	}

	return ret
}

// Convert ip string to uint32.
func IPv4toUint32(ip string) uint32 {
	var ret uint32 = 0
	for _, par := range strings.Split(ip, ".") {
		// shift the previously parsed bits over by one byte
		ret = ret << 8
		// set the low order bits to the current octet
		i, _ := strconv.Atoi(par)
		ret |= uint32(i)
	}

	return ret
}

// Convert int32 value to string.
func RuneToIPv4(ip int32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ip >> 24 & 0xff),
		(ip >> 16 & 0xff),
		(ip >> 8 & 0xff),
		(ip & 0xff))
}

// Convert uint32 value to string.
func Uint32ToIPv4(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ip >> 24 & 0xff),
		(ip >> 16 & 0xff),
		(ip >> 8 & 0xff),
		(ip & 0xff))
}

// min
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// max
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

//zero-garbage
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//zero-garbage
func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
