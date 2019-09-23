package xsf

import (
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	CallBackSessionTimeOut = &CallBackExceptionImpl{"xsf-server", "session timeout"}
)

type CallBackException interface {
	Caller() string
	Exception() string
}

type CallBackExceptionImpl struct {
	caller    string
	exception string
}

func (l *CallBackExceptionImpl) Caller() string {
	return l.caller
}

func (l *CallBackExceptionImpl) Exception() string {
	return l.exception
}

type Node struct {
	//每个session的数据信息
	DeadTime   time.Time
	Callback   func(sessionTag interface{}, svcData interface{}, exception ...CallBackException) //session的回调函数
	sessionTag interface{}
	Data       interface{}
}

func (c *Node) Task() { //执行callback的调度函数
	c.Callback(c.sessionTag, c.Data)
}

const (
	move = 0
	del  = 1
)

type taskItem struct {
	op   int
	node *list.Element
}
type taskDelay struct {
	f func()
}

//sessionManager的主要数据结构
type SessionManager struct {
	reportInterval  time.Duration //当策略为0即定时上报时，此为上报时间间隔 缺省1000ms，当策略为1即根据授权波动变化是，此值代表检查波动值的时间间隔
	rollTimeout     time.Duration //sessionManager内部遍历超时session的时间间隔
	Timeout         time.Duration //每个session的超时时间，超时后自动删除并触发回调函数，单位毫秒
	l               *list.List    //此链表用以存储节点之间先后顺序，队头入，队尾出
	lMu             *sync.Mutex
	m               map[interface{}]*list.Element //此map用以存储节点的位置信息，便于快速索引
	mMu             *sync.Mutex
	in              chan *taskItem   //队头写入管道，防止阻塞
	asyncDeleteTask chan *taskDelay  //延迟处理管道
	asyncUpdateTask chan *taskDelay  //延迟处理管道
	out             chan interface{} //队尾删除管道，防止阻塞
	bc              BootConfig
	MaxLic          int32 //总授权数
	//----------------------------------------------------------
	NowLic     int32          //当前的授权数
	BestLic    int32          //最佳的授权数
	reporter   *LbAdapter     //负载上报器
	reporterV2 *hermesAdapter //负载上报器
	logger     *utils.Logger
	strategy   int //可选值为0、1，分别代表定时上报、根据授权范围上报
	wave       int //波动值，当授权数变化值大于等于该值时，出发触发上报行为

	//-------------------------------------------------------
	taskChannelSize int //update 和 delete delay 管道的缓冲区大小,倍数
	taskSize        int //callback的消费协程数

	//-------------------------------------------------------
	logout bool //表示是否已经主动下线，若主动下线，则忽略后续的负载上报

	//--------------------------------------------------------------
	ctxCloseOut       context.Context //通知closeOut已经停止汇报
	ctxCloseOutCancel func()          //通知closeOut已经停止汇报
	ctxReport         context.Context //通知report停止上报
	ctxReportCancel   func()          //通知report停止上报
}

const (
	ReportOnTime          = 0
	ReportOnAuthorization = 1
	ReportOnIdles         = 2 //hermes
)

type SessionManagerOpt func(*SessionManager)

func WithSessionManagerBc(bc BootConfig) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.bc = bc
	}
}
func WithSessionManagerMaxLic(maxLic int32) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.MaxLic = maxLic
	}
}
func WithSessionManagerBestLic(bestLic int32) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.BestLic = bestLic
	}
}
func WithSessionManagerTimeout(timeout time.Duration) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.Timeout = timeout
	}
}
func WithSessionManagerRollTime(rollTime time.Duration) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.rollTimeout = rollTime
	}
}
func WithSessionManagerReportInterval(reportInterval int32) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.reportInterval = time.Duration(reportInterval) * time.Millisecond
	}
}
func WithSessionManagerReporter(reporter *LbAdapter) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.reporter = reporter
	}
}
func WithSessionManagerReporterv2(reporterv2 *hermesAdapter) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.reporterV2 = reporterv2
	}
}
func WithSessionManagerLogger(logHandle *utils.Logger) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.logger = logHandle
	}
}
func WithSessionManagerStrategy(strategy int) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.strategy = strategy
	}
}
func WithSessionManagerWave(wave int) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.wave = wave
	}
}

///////////////////////////////////////////////////////
func WithSessionManagerTaskSize(taskSize int) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.taskSize = taskSize
	}
}
func WithSessionManagerTaskChannelSize(taskChannelSize int) SessionManagerOpt {
	return func(sm *SessionManager) {
		sm.taskChannelSize = taskChannelSize
	}
}
func rstFormatter(rst []callRstItem) (interface{}) {
	var r []string
	for ix, rstItem := range rst {
		if rstItem.e == nil && rstItem.errcode == 0 {
			r = append(r, fmt.Sprintf("->%v:addr:%v success", ix, rstItem.addr))
			continue
		}

		r = append(r, fmt.Sprintf("->%v:addr:%v failed,s:%v,errcode:%v,e:%v",
			ix, rstItem.addr, func(in *Res) interface{} {
				if in == nil {
					return nil
				}
				return fmt.Sprintf("params:%v,data:%+v", in.GetAllParam(), in.GetData())
			}(rstItem.s), rstItem.errcode, rstItem.e))
	}
	return strings.Join(r, "||")
}

///////////////////////////////////////////////////////

//此处返回的error值为冗余值，可去掉
//NewSessionManager(bc BootConfig, maxLic int32, bestLic int32, timeout time.Duration, rollTime time.Duration, reportInterval time.Duration, reporter LbAdapter, logHandle *utils.Logger, strategy int, wave int)
func NewSessionManager(opts ...SessionManagerOpt) (*SessionManager, error) {
	var sm SessionManager
	callback, initErr := sm.Init(opts...)
	fcDelayInst.add(callback)
	addKillerCheck(killerFirstPriority, "sessionManager_"+sm.bc.CfgData.Service, &sm)
	return &sm, initErr
}

//sessionManager初始化函数
//timeout为session超时时间
//rollTime为遍历超时session的时间间隔
func (s *SessionManager) Init(opts ...SessionManagerOpt) (callback func() error, err error) {
	for _, opt := range opts {
		opt(s)
	}
	s.l = list.New()
	s.mMu = &sync.Mutex{}
	s.lMu = &sync.Mutex{}
	s.m = make(map[interface{}]*list.Element)
	s.in = make(chan *taskItem, s.MaxLic*10)
	s.asyncDeleteTask = make(chan *taskDelay, s.MaxLic*int32(s.taskChannelSize))
	s.asyncUpdateTask = make(chan *taskDelay, 1)
	s.out = make(chan interface{}, s.MaxLic*10)

	s.ctxCloseOut, s.ctxCloseOutCancel = context.WithCancel(context.Background())
	s.ctxReport, s.ctxReportCancel = context.WithCancel(context.Background())

	go s.dealTimeout()
	go s.inFront()
	go s.deleteElement()
	go s.delayDeleteWorker()
	go s.delayUpdateWorker()

	//引擎上线
	callback = func() error {
		s.logger.Infof("about to call fc delay task")

		/*
			老版本lb的上报
		*/
		if 0 != s.reporter.able {
			loginErr := s.reporter.Login(
				s.bc.CfgData.Service,
				s.MaxLic,
				atomic.LoadInt32(&s.MaxLic)-atomic.LoadInt32(&s.NowLic),
				s.BestLic,
				map[string]string{"bc": func() string {
					var m = make(map[string]interface{}, 10)
					m["CfgMode"] = s.bc.CfgMode
					m["CfgName"] = s.bc.CfgData.CfgName
					m["CfgDefault"] = s.bc.CfgData.CfgDefault
					m["Project"] = s.bc.CfgData.Project
					m["Group"] = s.bc.CfgData.Group
					m["Service"] = s.bc.CfgData.Service
					m["Version"] = s.bc.CfgData.Version
					m["CompanionUrl"] = s.bc.CfgData.CompanionUrl
					res, _ := json.Marshal(s.bc)
					return string(res)
				}()},
			)
			if nil != loginErr {
				log.Panic(fmt.Sprintf("s.reporter.Login fail"))
			}
		}

		/*
			如果两个lb都没有启用，就不要上报了
		*/
		if 0 != s.reporter.able || s.reporterV2.able {
			go s.report()
		}
		return nil
	}
	return
}

func (s *SessionManager) report() {
	s.logger.Debugf("SessionManager -> about to call report.")
	switch s.strategy {
	case ReportOnTime:
		{
			{
				idle := atomic.LoadInt32(&s.MaxLic) - atomic.LoadInt32(&s.NowLic)
				s.logger.Debugf("SessionManager -> s.update(s.MaxLic:%v, idle:%v, s.BestLic:%v)",
					s.MaxLic, idle, s.BestLic)
				_ = s.update(s.MaxLic, idle, s.BestLic)
			}

			ticker := time.NewTicker(s.reportInterval)
			for {
				select {
				case <-ticker.C:
					{
						idle := atomic.LoadInt32(&s.MaxLic) - atomic.LoadInt32(&s.NowLic)
						s.logger.Debugf("SessionManager -> s.update(s.MaxLic:%v, idle:%v, s.BestLic:%v)",
							s.MaxLic, idle, s.BestLic)
						_ = s.update(s.MaxLic, idle, s.BestLic)
					}
				}
			}
		}
	case ReportOnAuthorization:
		{
			lastMilepost := atomic.LoadInt32(&s.NowLic)
			ticker := time.NewTicker(s.reportInterval)
			for {
				select {
				case <-ticker.C:
					{
						NowLic := atomic.LoadInt32(&s.NowLic)
						if math.Abs(float64(NowLic-lastMilepost)) >= float64(s.wave) {
							idleInst := s.MaxLic - NowLic
							s.logger.Debugf("SessionManager -> lastMilepost:%v,s.update(s.MaxLic:%v, idleInst:%v, s.BestLic:%v)",
								lastMilepost, s.MaxLic, idleInst, s.BestLic)
							_ = s.update(s.MaxLic, idleInst, s.BestLic)
							lastMilepost = NowLic
						}
					}
				}
			}
		}
	case ReportOnIdles:
		{
			{
				err := s.updateV2(s.getAuthInfo)
				if err != nil {
					s.logger.Errorw("updateV2", "err", err)
				}
			}

			ticker := time.NewTicker(s.reportInterval)
		end:
			for {
				select {
				case <-ticker.C:
					{
						err := s.updateV2(s.getAuthInfo)
						if err != nil {
							s.logger.Errorw("updateV2", "err", err)
						}
					}
				case <-s.ctxReport.Done():
					{
						s.logger.Infow("closeOut notice stop reporting")
						s.ctxCloseOutCancel()
						break end
					}

				}
			}
		}
	default:
		{
			panic(fmt.Sprintf("the strategy->%v is illegal", s.strategy))
		}
	}
	s.logger.Infow("fn:sessionManager,report exit")

}
func (s *SessionManager) Closeout() {
	var closeOutCnt int64 = 0
	logId := fmt.Sprintf("closeout@%d%d", atomic.AddInt64(&closeOutCnt, 1), time.Now().Nanosecond())

	s.logger.Infow("notice report stop reporting", "logId", logId)
	s.ctxReportCancel()
	select {
	case <-s.ctxCloseOut.Done():
		{
			s.logger.Infow("already stop reporting", "logId", logId)
		}
	case <-time.After(time.Second * 3):
		{
			s.logger.Errorw("can't receive stop reporting signal", "logId", logId)
		}
	}

	if s.reporter.able != 0 || s.reporterV2.able {
		s.logger.Debugw("SessionManager -> about to call s.LoginOut()", "logId", logId)
		s.LoginOut(logId)
	}

	s.logger.Debugw("SessionManager -> about to deal lic", "logId", logId)

	//等待引擎授权释放
	ctx, _ := context.WithTimeout(context.Background(), s.Timeout*2)
	s.logger.Infow(
		"sessionManager check lic",
		"timeout", s.Timeout, "logId", logId)
end:
	for {
		select {
		case <-ctx.Done():
			{
				s.logger.Infow(
					"sessionManager timeout",
					"timeout", s.Timeout,
					"NowLic", atomic.LoadInt32(&s.NowLic), "MaxLic", atomic.LoadInt32(&s.MaxLic), "logId", logId)
				break end
			}
		default:
			{
				if atomic.LoadInt32(&s.NowLic) == 0 {
					s.logger.Infow("sessionManager exit",
						"timeout", s.Timeout,
						"NowLic", atomic.LoadInt32(&s.NowLic), "MaxLic", atomic.LoadInt32(&s.MaxLic), "logId", logId)
					break end
				}
				s.logger.Infow("sessionManager wait",
					"timeout", s.Timeout,
					"NowLic", atomic.LoadInt32(&s.NowLic), "MaxLic", atomic.LoadInt32(&s.MaxLic), "logId", logId)

				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}
func (s *SessionManager) getAuthInfo() (maxLic, idle, bestLic int32) {

	maxLic = s.MaxLic
	idle = atomic.LoadInt32(&s.MaxLic) - atomic.LoadInt32(&s.NowLic)
	bestLic = s.BestLic

	return
}
func (s *SessionManager) Update() error {
	if s.logout {
		return fmt.Errorf("logout called")
	}

	{
		///////////////////////////////////////////////////////////////
		switch s.strategy {
		case ReportOnTime, ReportOnAuthorization:
			{
				return s.update(s.MaxLic, atomic.LoadInt32(&s.MaxLic)-atomic.LoadInt32(&s.NowLic), s.BestLic)
			}
		case ReportOnIdles:
			{
				return s.updateV2(s.getAuthInfo)
			}
		default:
			{
				panic(fmt.Sprintf("the strategy->%v is illegal", s.strategy))
			}
		}

		///////////////////////////////////////////////////////////////
	}

}

func (s *SessionManager) UpdateAsync() {
	s.UpdateDelay()
}

func (s *SessionManager) UpdateDelay() {
	select {
	case s.asyncUpdateTask <- &taskDelay{f: func() {
		s.Update()
	}}:
		{
		}
	default:
		{
			s.logger.Debugf("throw update task away")
		}

	}
}
func (s *SessionManager) update(totalInst, idleInst, bestInst int32) error {
	s.logger.Debugf("SessionManager -> s.reporter.update(totalInst:%v, idleInst:%v, bestInst:%v)",
		totalInst, idleInst, bestInst)
	return s.reporter.Update(totalInst, idleInst, bestInst)
}

func (s *SessionManager) updateV2(getAuthInfo func() (int32, int32, int32)) error {

	return s.reporterV2.report(getAuthInfo)

}
func (s *SessionManager) LoginOut(logId string) (err error) {
	s.logger.Debugw(
		"SessionManager -> enter s.reporter.LoginOut()",
		"logId", logId)
	s.logout = true

	switch s.strategy {
	case ReportOnTime, ReportOnAuthorization:
		{
			s.logger.Debugw(
				"SessionManager -> s.reporter.LoginOut()",
				"strategy", "ReportOnTime、ReportOnAuthorization", "logId", logId)

			err = s.reporter.LoginOut()
		}
	case ReportOnIdles:
		{

			s.logger.Debugw(
				"SessionManager -> s.reporterV2.offline()",
				"strategy", "ReportOnIdles", "logId", logId)

			err := s.reporterV2.offline()
			if err != nil {
				err = fmt.Errorf("s.reporterV2.offline()->err:%v", err)
			}
		}
	default:
		{
			panic(fmt.Sprintf("the strategy->%v is illegal", s.strategy))
		}
	}

	s.logger.Debugw(
		"SessionManager -> leave s.reporter.LoginOut()",
		"strategy", "ReportOnIdles", "logId", logId)

	return
}

//通过sessionTag索引相关的sessiondata
func (s *SessionManager) GetSessionData(sessionTag interface{}) (interface{}, error) {
	s.mMu.Lock()
	sessionData, sessionDataOk := s.m[sessionTag]
	s.mMu.Unlock()

	if !sessionDataOk {
		return nil, fmt.Errorf("sessionData -> %v,sessionDataOk ->%v,can't get the value of sessionTag %s from map",
			sessionData, sessionDataOk, sessionTag)
	}
	sessionDataNode, sessionDataNodeOk := sessionData.Value.(*Node)

	if !sessionDataNodeOk {

		return nil, errors.New(fmt.Sprintf("can't Transform sessionData.Value -> %v Node.", sessionData.Value))

	} else {

		sessionDataNode.DeadTime = time.Now().Add(s.Timeout)

		select {
		case s.in <- &taskItem{op: move, node: sessionData}:
			{
				//s.logger.Debugf("write data to in success")
			}
		default:
			{
				//s.logger.Debugf("the queue in full")
			}
		}

	}
	return sessionDataNode.Data, nil
}

//向sessionManager中写入相关的session数据
func (s *SessionManager) SetSessionData(
	sessionTag interface{},
	svcData interface{},
	callback func(
	sessionTag interface{},
	svcData interface{},
	Exception ...CallBackException,
)) error {
	if atomic.AddInt32(&s.NowLic, 1) > atomic.LoadInt32(&s.MaxLic) {
		atomic.AddInt32(&s.NowLic, -1)
		return fmt.Errorf("lack of license -> now:%v,max:%v", s.NowLic, s.MaxLic)
	}
	{
		//debug
		s.logger.Debugf("op:%v,tag:%v,now:%v,deadTime:%v,timeout:%v", "SetSessionData",
			sessionTag, time.Now(), time.Now().Add(s.Timeout), s.Timeout)
	}
	node := &Node{
		Callback:   callback,
		Data:       svcData,
		sessionTag: sessionTag,
		DeadTime:   time.Now().Add(s.Timeout),
	}

	//取消异步机制采用同步机制
	s.mMu.Lock()
	_, sessionDataOk := s.m[node.sessionTag]
	s.mMu.Unlock()
	if sessionDataOk {
		atomic.AddInt32(&s.NowLic, -1)
		return fmt.Errorf("handle:%v is already stored", node.sessionTag)
	}
	s.mMu.Lock()
	s.lMu.Lock()
	s.m[node.sessionTag] = s.l.PushFront(node)
	s.lMu.Unlock()
	s.mMu.Unlock()
	return nil
}

//通过sessionTag删除相关的sessionData
func (s *SessionManager) DelSessionData(sessionTag interface{}, exception ...CallBackException) {
	s.logger.Infof("DelSessionData ->timestamp:%v,sessionTag:%v", time.Now(), sessionTag)
	////////////////////////////////////////////////
	/*
		1、清空对应sessionTag的map
		2、清空对应sessionTag的list
		3、执行sessionTag对应的callback
	*/
	s.mMu.Lock()
	elem, elemOk := s.m[sessionTag]
	if !elemOk {
		s.mMu.Unlock()
		return
	}
	delete(s.m, sessionTag)
	s.mMu.Unlock()
	/*删除链表*/
	//s.mMu.Lock()
	//s.l.Remove(elem)
	s.in <- &taskItem{op: del, node: elem}
	/*删除map*/
	//s.mMu.Lock()
	//s.mMu.Unlock()

	sessionDataNode, _ := elem.Value.(*Node)
	sessionDataNode.Callback(sessionDataNode.sessionTag, sessionDataNode.Data, exception...)

	atomic.AddInt32(&s.NowLic, -1)
}

func (s *SessionManager) DelSessionDataDelay(sessionTag interface{}, exception ...CallBackException) {
	s.asyncDeleteTask <- &taskDelay{f: func() {
		s.logger.Infof("DelSessionDataDelay ->timestamp:%v,sessionTag:%v", time.Now(), sessionTag)
		////////////////////////////////////////////////
		/*
			1、清空对应sessionTag的map
			2、清空对应sessionTag的list
			3、执行sessionTag对应的callback
		*/
		s.mMu.Lock()
		elem, elemOk := s.m[sessionTag]
		if !elemOk {
			s.mMu.Unlock()
			return
		}
		delete(s.m, sessionTag)
		s.mMu.Unlock()
		/*删除链表*/
		//s.mMu.Lock()
		//s.l.Remove(elem)
		s.in <- &taskItem{op: del, node: elem}
		/*删除map*/
		//s.mMu.Lock()
		//s.mMu.Unlock()

		sessionDataNode, _ := elem.Value.(*Node)
		sessionDataNode.Callback(sessionDataNode.sessionTag, sessionDataNode.Data, exception...)

		atomic.AddInt32(&s.NowLic, -1)
	}}
}

func (s *SessionManager) delayDeleteWorker() {
	s.logger.Debugf("SessionManager -> about to call inFront()")
	s.logger.Errorf("taskSize:%v,taskChannelSize:%v\n", s.taskSize, s.taskChannelSize)
	// 队头入
	var wg sync.WaitGroup
	for gCnt := 0; gCnt < s.taskSize; gCnt++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case taskDelete := <-s.asyncDeleteTask:
					{
						s.logger.Infof("SessionManager taskDelete -> taskUpdate.f()")
						taskDelete.f()
					}
				}
			}
		}()
	}
	wg.Wait()

}
func (s *SessionManager) delayUpdateWorker() {
	s.logger.Debugf("SessionManager -> about to call inFront()")
	// 队头入
	for {
		select {
		case taskUpdate := <-s.asyncUpdateTask:
			{
				s.logger.Infof("SessionManager taskUpdate -> taskUpdate.f()")
				taskUpdate.f()
			}
		}
	}
}

//队头维护函数，将数据提交到链表头
func (s *SessionManager) inFront() {
	s.logger.Debugf("SessionManager -> about to call inFront()")
	// 队头入
	for {
		select {
		case task := <-s.in:
			{
				switch task.op {
				case move:
					{
						s.lMu.Lock()
						s.l.MoveToFront(task.node)
						s.lMu.Unlock()
					}
				case del:
					{
						s.lMu.Lock()
						s.l.Remove(task.node)
						s.lMu.Unlock()
					}
				}
			}
		}
	}
}

//队列维护函数，根据sessionTag删除sessionManager中的相关数据
func (s *SessionManager) deleteElement() {
	s.logger.Debugf("SessionManager -> about to call deleteElement()")
	// 队尾出
	for {
		select {
		case sessionNode := <-s.out:
			{
				node, nodeOk := sessionNode.(*Node)
				if !nodeOk {
					s.logger.Errorf("can't convert sessionNode to *Node")
				}
				s.logger.Infow("begin call node.callback", "sessionTag",
					node.sessionTag, "fn", "deleteElement")
				node.Callback(node.sessionTag, node.Data, CallBackSessionTimeOut)
				s.logger.Infow("end call node.callback", "sessionTag",
					node.sessionTag, "fn", "deleteElement")
				atomic.AddInt32(&s.NowLic, -1)
			}
		}
	}
}

//超时维护函数，将超时session提交给delete element函数
func (s *SessionManager) dealTimeout() {
	s.logger.Infof("SessionManager rollTimeout:%v -> about to call dealTimeout()",
		s.rollTimeout)
	timer := time.NewTicker(s.rollTimeout)
	for {
		select {
		case <-timer.C:
			{
				s.dealTimeWorker()
			}
		}
	}
}

func (s *SessionManager) dealTimeWorker() {
	s.logger.Infof("SessionManager -> in dealTimeout.pos:timer.C")
	ts := time.Now()
	for {
		s.lMu.Lock()
		element := s.l.Back()

		if element == nil {
			s.lMu.Unlock()
			break
		}
		sessionDataNode, _ := element.Value.(*Node)

		if ts.After(sessionDataNode.DeadTime) {
			s.logger.Infof("dealTimeout ->timestamp:%v,sessionTag:%v",
				ts, sessionDataNode.sessionTag)
			s.l.Remove(element)
			s.lMu.Unlock()
			s.mMu.Lock()
			if _, ok := s.m[sessionDataNode.sessionTag]; !ok {
				s.logger.Errorf("can't take of tag:%v of sessionDataNode",
					sessionDataNode.sessionTag)
				s.mMu.Unlock()

			} else {
				delete(s.m, sessionDataNode.sessionTag)
				s.mMu.Unlock()
				s.out <- sessionDataNode
			}

		} else {
			s.lMu.Unlock()
			break
		}
	}
}
