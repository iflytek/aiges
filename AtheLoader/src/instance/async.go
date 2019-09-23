package instance

import (
	"errors"
	"sync"
)

var (
	errRepeat  = errors.New("wrapper handle repeat acquire async chan")
	errInvalid = errors.New("invalid wrapper handle, can't find channel")
)

var wrapperAsyncChan map[string]chan ActMsg
var wacMutex sync.Mutex

func init() {
	wrapperAsyncChan = make(map[string]chan ActMsg)
}

func AllocChan(hdl string) (amc chan ActMsg, err error) {
	wacMutex.Lock()
	defer wacMutex.Unlock()
	_, exist := wrapperAsyncChan[hdl]
	if exist {
		return nil, errRepeat
	}
	ac := make(chan ActMsg, seqRltSize)
	wrapperAsyncChan[hdl] = ac
	return ac, nil
}

func FreeChan(hdl string) (err error) {
	wacMutex.Lock()
	defer wacMutex.Unlock()
	ac, exist := wrapperAsyncChan[hdl]
	if exist {
		close(ac)
		delete(wrapperAsyncChan, hdl)
	}
	return
}

func QueryChan(hdl string) (amc chan ActMsg, err error) {
	wacMutex.Lock()
	defer wacMutex.Unlock()
	ac, exist := wrapperAsyncChan[hdl]
	if !exist {
		return nil, errInvalid
	}
	return ac, nil
}
