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
	"errors"
	"github.com/xfyun/flume"
	"github.com/xfyun/thrift"
	"net"
	"reflect"
	"time"
)

type flumeClient struct {
	// flume agent host
	host string
	// flume agent port
	port string
	// socket
	transport *thrift.TSocket
	// flume sdk client
	rpcClient *flume.ThriftSourceProtocolClient
}

func (fc *flumeClient) open() bool {
	// check have prepared
	if fc.rpcClient == nil || fc.transport == nil {
		if err := fc.prepare(); err != nil {
			return false
		}
	}

	// check is open
	if fc.transport.IsOpen() {
		return true
	}

	// reopen
	if err := fc.transport.Open(); err != nil {
		errorf("thrift transport open error %v.", err)
		if reflect.TypeOf(err).String() == "*thrift.tTransportException" {
			if err.(thrift.TTransportException).TypeId() == thrift.ALREADY_OPEN {
				return true
			} else {
				errorf("thrift transport error %v.", err)
				return false
			}
		}
	}

	return false
}

func (fc *flumeClient) close() bool {
	if err := fc.transport.Close(); err != nil {
		return false
	}
	return true
}

// prepare flume transport
func (fc *flumeClient) prepare() error {
	// transport settings, see java syringe
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTCompactProtocolFactory()

	var err error
	fc.transport, err = thrift.NewTSocket(net.JoinHostPort(fc.host, fc.port))
	if err != nil {
		errorf("thrift transport create error : %v.", err)
		return errors.New("thrift transport build err")
	}

	fc.transport.SetTimeout(time.Second * 10)
	useTransport := transportFactory.GetTransport(fc.transport)
	if useTransport == nil {
		return errors.New("thrift not get transport")
	}

	fc.rpcClient = flume.NewThriftSourceProtocolClientFactory(useTransport, protocolFactory)
	if fc.rpcClient == nil {
		return errors.New("thrift new client failed")
	}

	return nil
}
