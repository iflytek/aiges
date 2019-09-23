package xsf

import (
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/VividCortex/ewma"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type vCpuBucket struct {
	sum   int64
	count int64
}
type vCpuWin struct {
	timeSlices []vCpuBucket

	timeSliceSize int64

	timePerSlice time.Duration

	winSize int64

	cursor int64

	preTs time.Time

	winDur time.Duration
}

func getVCpuCfg(
	cfgVer string,
	cfgPrj string,
	cfgGroup string,
	cfgService string,
	cfgName string,
	cfgUrl string,
	cfgMode utils.CfgMode) (vCpuMap map[string]interface{}, err error) {
	logCfgOpt := &utils.CfgOption{}
	utils.WithCfgVersion(cfgVer)(logCfgOpt)
	utils.WithCfgPrj(cfgPrj)(logCfgOpt)
	utils.WithCfgGroup(cfgGroup)(logCfgOpt)
	utils.WithCfgService(cfgService)(logCfgOpt)
	utils.WithCfgName(cfgName)(logCfgOpt)
	utils.WithCfgURL(cfgUrl)(logCfgOpt)
	cfg, err := utils.NewCfg(cfgMode, logCfgOpt)

	if err != nil {
		return nil, err
	}
	cpuSec, ok := cfg.GetSection(VCPUSEC).(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't convert %#v to map[string]interface{]", cfg.GetSection("cpu"))
	}
	return cpuSec, nil
}
func newVCpuWindow(timePerSlice time.Duration, winSize int64) *vCpuWin {
	vCpuWindowInst := vCpuWin{}
	vCpuWindowInst.Init(timePerSlice, winSize)
	return &vCpuWindowInst
}

func (s *vCpuWin) reset() {

	ticker := time.NewTicker(s.winDur)

	for {
		select {
		case <-ticker.C:
			{
				if time.Now().Sub(s.preTs) < s.winDur {
					continue
				}
				for index := 0; index < len(s.timeSlices); index++ {
					atomic.StoreInt64(&s.timeSlices[index].sum, 0)
				}
			}
		}
	}
}
func (s *vCpuWin) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]vCpuBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *vCpuWin) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

func (s *vCpuWin) addSampling(sampling int64) {
	var index = s.locationIndex()
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		atomic.AddInt64(&s.timeSlices[index].sum, sampling)
		atomic.AddInt64(&s.timeSlices[index].count, 1)
	} else {
		atomic.AddInt64(&s.timeSlices[index].sum, sampling)
		atomic.AddInt64(&s.timeSlices[index].count, 1)

		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

func (s *vCpuWin) getStats() []float64 {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].count, 0)

		s.clearBetween(oldCursor, index)
	}

	var avgs []float64
	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
		avgs = append(avgs, float64(atomic.LoadInt64(&bucketTmp.sum))/float64(atomic.LoadInt64(&bucketTmp.count)))
	}

	return avgs
}

func (s *vCpuWin) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].count, 0)
	}
}

const remainVCpuChanSize = 1

type VCpuManager struct {
	virtualCpuRwMu sync.RWMutex
	virtualCpuPair map[string]int64 //型号+虚拟cpu

	cpuModel string //当前CPU型号
	cpuScore int64  //当前CPU满血虚拟核心，暂定从映射表里获取，如果映射表出错，采用物理核心替代

	freq              float64      //cpu主频
	totalPhysicalCpus int64        //总的物理核心数
	physicalCpus      float64      //被分配的cpu物理核数
	virtualCpus       float64      //被分配的虚拟CPU，由物理cpu核数转换而来
	remainVCpu        float64      //剩余的虚拟CPU,经过拟合后的值
	RemainVCpuInt64   int64        ////由remainVCpu四舍五入而来
	remainVCpuChan    chan float64 //通知主线程过来取

	//window *vCpuWin //数据窗口,2019-08-01 15:43:12 移除数据窗口

	interval          time.Duration
	hardWareCollector *utils.HardWareCollector

	fittingManager ewma.MovingAverage
}

//func NewVCpuManager(interval time.Duration, winSize int64, vCpuIn map[string]interface{}) (*VCpuManager, error) {
//	vCpuInst := VCpuManager{}
//	return &vCpuInst, vCpuInst.Init(interval, winSize, vCpuIn)
//}
func NewVCpuManager(interval time.Duration, vCpuIn map[string]interface{}) (*VCpuManager, error) {
	vCpuInst := VCpuManager{}
	return &vCpuInst, vCpuInst.Init(interval, vCpuIn)
}

func (v *VCpuManager) Init(interval time.Duration, vCpuIn map[string]interface{}) error {
	{
		err := v.Store(vCpuIn)
		if nil != err {
			return err
		}
	}
	{
		//拟合初始化
		v.fittingManager = ewma.NewMovingAverage(VCPUFITTEDVALUE)
	}
	{
		var err error
		v.hardWareCollector, err = utils.NewHardWareCollector()
		if nil != err {
			return err
		}
		//v.window = newVCpuWindow(interval, winSize)
		v.interval = interval
	}
	{
		//主频
		MHZs, MHZsErr := v.hardWareCollector.GetCpuMHz()
		if nil != MHZsErr {
			return MHZsErr
		}
		if 0 == len(MHZs) {
			return fmt.Errorf("GetCpuMHz return is empty")
		}
		v.freq = MHZs[0]
		loggerStd.Printf("cpuMHz:%v\n", v.freq)
	}
	{
		//计算物理总核数
		physicalCpus, physicalCpusErr := v.hardWareCollector.GetCpus()
		if nil != physicalCpusErr {
			return fmt.Errorf("GetCpus return is err:%v", physicalCpusErr)
		}
		v.totalPhysicalCpus = int64(physicalCpus)
	}
	{
		//被分配的核数
		c, e := v.hardWareCollector.GetCpuLimit()
		if nil != e {
			return fmt.Errorf("GetCpuLimit return is err:%v", e)
		}
		if c < 0 {
			v.physicalCpus = float64(v.totalPhysicalCpus)
		} else {
			v.physicalCpus = c
		}
		loggerStd.Printf("cpulimit:%v\n", v.physicalCpus)
	}
	{
		//计算虚拟cpu
		Models, ModelsErr := v.hardWareCollector.GetModel()
		if nil != ModelsErr {
			return ModelsErr
		}
		if 0 == len(Models) {
			return fmt.Errorf("GetModel return is empty")
		}
		v.cpuModel = Models[0]
		loggerStd.Printf("cpuModel:%v\n", v.cpuModel)

		v.cpuScore = v.LoadWithModel(v.cpuModel)
		dbgLoggerStd.Printf("from LoadWithModel,cpuScore:%v\n", v.cpuScore)
		if 0 == v.cpuScore {
			dbgLoggerStd.Printf("cpuScore==0,use physicalCpus,cpuScore:%v\n", v.cpuScore)
			v.cpuScore = v.totalPhysicalCpus
		}
		loggerStd.Printf("cpuScore:%v\n", v.cpuScore)
	}
	{
		//计算被分配的虚拟CPU
		if 0 != v.cpuScore {
			v.virtualCpus = (v.physicalCpus / float64(v.totalPhysicalCpus)) * float64(v.cpuScore)
		} else {
			//如果没有相关数据，则直接采用物理核
			v.virtualCpus = v.physicalCpus
		}
	}

	v.remainVCpuChan = make(chan float64, remainVCpuChanSize)

	go v.start(v.remainVCpuChan)
	go v.dataProcessing(v.remainVCpuChan)

	return nil
}
func (v *VCpuManager) fitting(dataIn float64) float64 {
	v.fittingManager.Add(dataIn)
	return v.fittingManager.Value()
}
func (v *VCpuManager) getRemainVCpu() (float64, error) {
	//获取当前的剩余的虚拟cpu数
	usage, usageErr := v.hardWareCollector.GetCpuUsage()
	if usageErr != nil {
		return 0, usageErr
	}
	usageRate := usage / float64(v.totalPhysicalCpus) / 100

	remainVCpu := v.virtualCpus - (float64(v.cpuScore) * usageRate)
	dbgLoggerStd.Printf("fn:getRemainVCpu,remainVCpu:%v,virtualCpus:%v,cpuScore:%v,usage:%v,usageRate:%v\n",
		remainVCpu, v.virtualCpus, v.cpuScore, usage, usageRate)

	return v.fitting(remainVCpu), nil
}

//避免节点接入过快
func (v *VCpuManager) setThresholdRate(vCpuIn float64) float64 {
	dur := time.Since(globalStart)
	rate := func() float64 {
		//todo 后序完善
		if dur > 10*time.Second {
			return 1
		} else if dur > 9*time.Second {
			return 0.9
		} else if dur > 8*time.Second {
			return 0.8
		} else if dur > 7*time.Second {
			return 0.7
		} else if dur > 6*time.Second {
			return 0.6
		} else if dur > 5*time.Second {
			return 0.5
		} else if dur > 4*time.Second {
			return 0.4
		} else if dur > 3*time.Second {
			return 0.3
		} else if dur > 2*time.Second {
			return 0.2
		} else if dur > 1*time.Second {
			return 0.1
		} else {
			return 0.05
		}
	}()
	dbgLoggerStd.Printf("fn:setThresholdRate,dur:%v,rate:%v,vCpu:%v\n", dur, rate, vCpuIn)
	return vCpuIn * rate
}

func (v *VCpuManager) dataProcessing(ch <-chan float64) {
	round := func(x float64) int64 {
		return int64(math.Floor(x + 0.5))
	}
	for {
		select {
		case data := <-ch:
			{
				atomic.StoreInt64(&(v.RemainVCpuInt64), round(v.setThresholdRate(data)))
				dbgLoggerStd.Printf("fn:dataProcessing,data_float64:%v,data_int64:%v\n",
					data, atomic.LoadInt64(&(v.RemainVCpuInt64)))
			}
		}
	}
}
func (v *VCpuManager) start(ch chan float64) {

	dbgLoggerStd.Printf("fn:start,interval:%v\n", v.interval)
	tickerInst := time.NewTicker(v.interval)

	for {
		dbgLoggerStd.Println("repeat VCpuManager start....")
		select {
		case <-tickerInst.C:
			{
				//通道的方式通知接收线程
				dbgLoggerStd.Println("begin to collect remain...")
				remainVCpu, remainVCpuErr := v.getRemainVCpu()
				dbgLoggerStd.Printf("fn:start,remainVCpu:%v,remainVCpuErr:%v\n",
					remainVCpu, remainVCpuErr)
				select {
				case ch <- remainVCpu:
					{
						dbgLoggerStd.Println("write data to ch normally...")
					}
				default:
					{
						dbgLoggerStd.Println("nobody consume...")
						//先消费再写入
					empty:
						for {
							select {
							case <-ch:
							default:
								break empty
							}
						}
						dbgLoggerStd.Println("write data to ch...")
						ch <- remainVCpu
					}
				}
			}
		}
	}
}
func (v *VCpuManager) Store(vCpuIn map[string]interface{}) error {
	dbgLoggerStd.Printf("fn:Store,vCpuIn:%+v\n", vCpuIn)
	v.virtualCpuRwMu.Lock()
	defer v.virtualCpuRwMu.Unlock()

	vCpuInst := make(map[string]int64)
	for k, v := range vCpuIn {
		vI, vIOk := v.(int64)
		if !vIOk {
			return fmt.Errorf("can't convert %#v to int64", v)
		}
		vCpuInst[k] = vI
	}
	v.virtualCpuPair = vCpuInst
	dbgLoggerStd.Printf("fn:Store,virtualCpuPair:%+v\n", v.virtualCpuPair)
	return nil
}
func (v *VCpuManager) LoadWithModel(cpuModel string) int64 {
	v.virtualCpuRwMu.RLock()
	defer v.virtualCpuRwMu.RUnlock()
	dbgLoggerStd.Printf("fn:LoadWithModel,virtualCpuPair:%+v\n", v.virtualCpuPair)
	return v.virtualCpuPair[cpuModel]
}
