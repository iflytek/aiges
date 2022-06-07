package xsf

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
	"github.com/xfyun/xsf/utils"
)

type QpsLimiter struct {
	reportInterval  time.Duration //上报时间，单位毫秒
	bc              BootConfig    //启动配置
	maxReqCount     int32         //最大请求数
	bestReqCount    int32         //最佳请求数
	currentReqCount int32         //当前的请求数
	reporter        LbAdapter     //负载上报器
	reporterv2      hermesAdapter //负载上报器
	logger          *utils.Logger
	interval        time.Duration //重置时间间隔，单位毫秒
	ticker          *time.Ticker  //根据interval创建的定时器
}
type QpsLimiterOpt func(*QpsLimiter)

func WithQpsLimiterReportInterval(reportInterval time.Duration) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.reportInterval = reportInterval
	}
}
func WithQpsLimiterBc(bc BootConfig) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.bc = bc
	}
}
func WithQpsLimiterMaxReqCount(maxReqCount int32) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.maxReqCount = maxReqCount
	}
}
func WithQpsLimiterBestReqCount(bestReqCount int32) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.bestReqCount = bestReqCount
	}
}
func WithQpsLimiterCurrentReqCount(currentReqCount int32) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.currentReqCount = currentReqCount
	}
}
func WithQpsLimiterReporter(reporter LbAdapter) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.reporter = reporter
	}
}
func WithQpsLimiterReporterv2(reporterv2 hermesAdapter) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.reporterv2 = reporterv2
	}
}
func WithQpsLimiterLogger(logger *utils.Logger) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.logger = logger
	}
}
func WithQpsLimiterInterval(interval int32) QpsLimiterOpt {
	return func(ql *QpsLimiter) {
		ql.interval = time.Duration(interval) * time.Millisecond
	}
}

//NewQpsLimiter(bc BootConfig, maxReq int32, bestReq int32, interval time.Duration, reportInterval time.Duration, reporter LbAdapter, logHandle *utils.Logger)
func NewQpsLimiter(opts ...QpsLimiterOpt) (*QpsLimiter, error) {
	var qpsLimiter QpsLimiter
	if SetQpsLimiterErr := qpsLimiter.SetQpsLimiter(opts...); SetQpsLimiterErr != nil {
		return nil, fmt.Errorf("SetQpsLimiter failed -> SetQpsLimiterErr:%v", SetQpsLimiterErr)
	}
	addKillerCheck(killerFirstPriority, "qpsLimiter_"+qpsLimiter.bc.CfgData.Service, &qpsLimiter)
	return &qpsLimiter, nil
}

//set limiter params
func (l *QpsLimiter) SetQpsLimiter(opts ...QpsLimiterOpt) error {
	for _, opt := range opts {
		opt(l)
	}
	if l.interval < 1 {
		return fmt.Errorf("interval -> %v is less than 1", l.interval)
	}
	l.ticker = time.NewTicker(l.interval)
	go l.guarder()
	return l.reporter.Login(l.bc.CfgData.Service, l.maxReqCount, l.maxReqCount-l.currentReqCount, l.bestReqCount, map[string]string{"bc": func() string { res, _ := json.Marshal(l.bc); return string(res) }()})
}

func (l *QpsLimiter) report() {
	l.logger.Debugf("QpsLimiter -> about to start report")
	ticker := time.NewTicker(l.reportInterval)
	for {
		select {
		case <-ticker.C:
			{
				idleInst := l.maxReqCount - l.currentReqCount
				l.logger.Debugf("about to call l.Update(l.maxReqCount:%v, idleInst:%v, l.bestReqCount:%v)",
					l.maxReqCount, idleInst, l.bestReqCount)
				l.Update(l.maxReqCount, idleInst, l.bestReqCount)
			}
		}
	}
}
func (l *QpsLimiter) Closeout() {
	l.logger.Debugf("QpsLimiter -> l.reporter.LoginOut()")
	l.reporter.LoginOut()
}
func (l *QpsLimiter) Update(maxReq, idleInst, bestReq int32) error {
	l.logger.Debugf("QpsLimiter -> l.reporter.Update(maxReq:%v, idleInst:%v, bestReq:%v)", maxReq, idleInst, bestReq)
	return l.reporter.Update(maxReq, idleInst, bestReq)
}
func (l *QpsLimiter) LoginOut() error {
	l.logger.Debugf("QpsLimiter -> l.reporter.LoginOut()")
	return l.reporter.LoginOut()
}

//start the limiter
func (l *QpsLimiter) guarder() {
	l.logger.Debugf("QpsLimiter -> about to start guarder")
	for {
		select {
		case <-l.ticker.C:
			{
				atomic.StoreInt32(&l.currentReqCount, 0)
			}
		}
	}
}

//check the qps
func (l *QpsLimiter) CheckQps() bool {
	if atomic.AddInt32(&l.currentReqCount, 1) < atomic.LoadInt32(&l.maxReqCount) {
		return true
	} else {
		return false
	}
}
