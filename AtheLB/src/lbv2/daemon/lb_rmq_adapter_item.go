package daemon

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"log"
	"strings"
)

const (
	nodeNotStarted = "No connection could be made because the target machine actively refused it"
	nodeRestarted  = "An established connection was aborted by the software in your host machine"
	nodeEOF        = "EOF"
	nodeNotOpen    = "Connection not open"
	nodeNoData     = "No more data"
	nodeListEmpty  = "clientList is empty"
)

type rmqAdapterItem struct {
	addr   string
	client *MTRMessageServiceClient
}

func (r *rmqAdapterItem) Init(addr string) (err error) {
	r.addr = addr
	transport, err := thrift.NewTSocket(addr)
	if nil != err {
		log.Fatal(err)
	}
	useTransport := thrift.NewTBufferedTransportFactory(100000).GetTransport(transport)
	r.client = NewMTRMessageServiceClientFactory(useTransport, thrift.NewTBinaryProtocolFactoryDefault())
	err = transport.Open()

	return
}

func (r *rmqAdapterItem) Produce(topic string, body string) (produceReply int64, produceErr error) {
	var msg MTRMessage
	msg.Topic = topic
	msg.Body = []byte(body)
	msg.Protocol = MTRProtocol_PERSONALIZED
	produceReply, produceErr = r.client.Produce(&msg, true)
	return
}

func (r *rmqAdapterItem) Consume(topic, group string) (consumeReply *MTRMessage, consumeErr error) {
	consumeReply, consumeErr = r.client.Consume(topic, group)
	if nil != consumeErr {
		resetErr := r.Reset(consumeErr)
		consumeErr = fmt.Errorf("addr:%s,rmqErr:%v,rmqReset:%v", r.addr, consumeErr, resetErr)
	}
	return
}
func (r *rmqAdapterItem) Finit() {
	_ = r.client.Transport.Close()
}
func (r *rmqAdapterItem) Reset(err error) error {
	if nil == err {
		return nil
	}
	if strings.Contains(err.Error(), nodeNotStarted) {
		r.Finit()
		return r.Init(r.addr)

	}
	if strings.Contains(err.Error(), nodeRestarted) {
		r.Finit()
		return r.Init(r.addr)

	}
	if strings.Contains(err.Error(), nodeEOF) {
		r.Finit()
		return r.Init(r.addr)
	}
	if strings.Contains(err.Error(), nodeNotOpen) {
		r.Finit()
		return r.Init(r.addr)
	}
	return nil
}
