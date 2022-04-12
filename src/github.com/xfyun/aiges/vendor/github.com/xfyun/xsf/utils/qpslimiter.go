/*
 @brief		this is a limiter for qps
 @file		qpslimiter
 @author	sqjain
 @version	1.0
 @date		2017.11.17
*/

// this is a limiter for qps.
package utils

import (
	"sync/atomic"
	"time"
	"errors"
)

type QpsLimiter struct {
	maxReqCount     int64
	currentReqCount int64
	ticker          *time.Ticker
}

func NewQpsLimiter(maxReq int64, interval int64) (*QpsLimiter, error) {
	var qpslimiter QpsLimiter
	if !qpslimiter.SetQpsLimiter(maxReq, interval) {
		return nil, errors.New("SetQpsLimiter failed")
	}
	return &qpslimiter, nil
}

//set limiter params
func (l *QpsLimiter) SetQpsLimiter(maxReq int64, interval int64) bool {
	if interval < 1 {
		return false
	}
	l.ticker = time.NewTicker(time.Duration(interval) * time.Second)
	atomic.StoreInt64(&l.maxReqCount, maxReq)
	go l.guarder()
	return true
}

//start the limiter
func (l *QpsLimiter) guarder() {
	for {
		select {
		case <-l.ticker.C:
			{
				atomic.StoreInt64(&l.currentReqCount, 0)
			}
		}
	}
}

//check the qps
func (l *QpsLimiter) CheckQps() bool {
	if atomic.AddInt64(&l.currentReqCount, 1) < atomic.LoadInt64(&l.maxReqCount) {
		return true
	} else {
		return false
	}
}
