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

import "fmt"

// custom log interface
type CustomLogInterface interface {
	// Infof formats message according to format specifier
	// and writes to log with level = Info.
	Infof(format string, params ...interface{})

	// Debugf formats message according to format specifier
	// and writes to log with level = Debug.
	Debugf(format string, params ...interface{})

	// Errorf formats message according to format specifier
	// and writes to log with level = Error.
	Errorf(format string, params ...interface{})
}

var (
	// logger
	Logger CustomLogInterface
)

// infof output a formatted info log
func infof(format string, params ...interface{}) {
	if Logger != nil {
		Logger.Infof("flange: "+format, params...)
	}
}

// debugf output a formatted debug log
func debugf(format string, params ...interface{}) {
	if Logger != nil {
		Logger.Debugf("flange: "+format, params...)
	}
}

// errorf output a formatted error log
func errorf(format string, params ...interface{}) {
	runtimeErrorInterpolation = fmt.Sprintf(format, params...)

	if Logger != nil {
		Logger.Errorf("flange: "+format, params...)
	}
}

type FmtLog struct {
}

func (log *FmtLog) Infof(format string, params ...interface{}) {
	fmt.Println(fmt.Sprintf(format, params...))
}

func (log *FmtLog) Debugf(format string, params ...interface{}) {
	fmt.Println(fmt.Sprintf(format, params...))
}

func (log *FmtLog) Errorf(format string, params ...interface{}) {
	fmt.Println(fmt.Sprintf(format, params...))
}
