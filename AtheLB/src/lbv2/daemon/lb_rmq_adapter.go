package daemon

import (
	"fmt"
	"strings"
	"sync/atomic"
)

var RmqAdapter rmqAdapter

type rmqAdapter struct {
	clientList []rmqAdapterItem
	ix         int64
}

var addrList []string

func (r *rmqAdapter) String() string {
	var addrs []string
	for _, v := range r.clientList {
		addrs = append(addrs, v.addr)
	}
	return strings.Join(addrs, ",")
}
func (r *rmqAdapter) Init(addr []string) (err error) {
	std.Println("rmqTargets:", addr)
	var errs []string

	for _, addrItem := range addr {
		addrList = append(addrList, addrItem)
		var rmqAdapterItemTmp rmqAdapterItem
		if err = rmqAdapterItemTmp.Init(addrItem); nil != err {
			errs = append(errs, fmt.Sprintf("addr:%v,err:%v", addrItem, err))
			r.clientList = append(r.clientList, rmqAdapterItemTmp) //此处基于代码可读性考虑，保留冗余代码
		} else {
			r.clientList = append(r.clientList, rmqAdapterItemTmp)
		}
	}
	if 0 == len(errs) {
		return nil
	}

	return fmt.Errorf("%v", strings.Join(errs, ","))
}

func (r *rmqAdapter) Produce(topic string, body string) (produceReply int64, produceErr error) {
	for _, rmqAdapterItemTmp := range r.clientList {
		produceReply, produceErr = rmqAdapterItemTmp.Produce(topic, body)
		if 0 == produceReply && nil == produceErr {
			break
		} else {
			continue
		}
	}
	return
}
func (r *rmqAdapter) Consume(topic, group string) (RemoteRmqAddr string, ConsumeR *MTRMessage, ConsumeE error) {
	var errors []string
	var flag = false
	for i := 0; i < len(r.clientList); i++ {
		consumeR, consumeE := r.clientList[atomic.AddInt64(&r.ix, 1)%int64(len(r.clientList))].Consume(topic, group)
		RemoteRmqAddr = r.clientList[atomic.LoadInt64(&r.ix)%int64(len(r.clientList))].addr
		if nil == consumeE {
			flag = true
			ConsumeR = consumeR
			break
		} else {
			errors = append(errors, fmt.Sprintf("%v\n", consumeE))
			continue
		}
	}
	if 0 == len(r.clientList) {
		ConsumeE = fmt.Errorf(nodeListEmpty)
	} else if len(errors) > 0 && !flag {
		ConsumeE = fmt.Errorf("%s", strings.Join(errors, ";"))
	}
	return
}
