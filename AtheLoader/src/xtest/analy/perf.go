package analy

import (
	"time"
)

/*
用于性能统计, 统计维度如下：
1. 首结果平均|最大|最小性能
2. 尾结果平均|最大|最小性能
3. 会话平均|最大|最小耗时
4. 各接口方法(AIIn/AIOut/AIExcp)平均|最大|最小耗时
以上数据提供分布区间统计结果
*/

// 计时类型,用于控制计时开关
const (
	FIFOPERF = 1 << iota // 首结果耗时
	LILOPERF             // 尾结果耗时
	SESSPERF             // 会话耗时
	INTFPERF             // 接口耗时
)

// 计时定点,用于标记计时位置
const (
	pointCreate int = 1 << iota
	pointUpBegin
	pointDownBegin
	pointUpEnd
	pointDownEnd
	pointDestroy
)

type PerfDetail struct {
	cTime     time.Time // create 时间
	dTime     time.Time // destroy 时间
	firstUp   time.Time
	lastUp    time.Time
	firstDown time.Time
	lastDown  time.Time
	upCost    []time.Time // 上行接口耗时
	downCost  []time.Time // 下行接口耗时

}

type perfDist struct {
	level   int                            // 性能统计等级, 最高：FIFOPERF | LILOPERF | SESSPERF | INTFPERF
	details map[string] /*sid*/ PerfDetail // 分布数据需要保存全量会话数据

}

func (pc *perfDist) Start(perfLevel int) {

}

func (pc *perfDist) TickPoint(point int) {
	// TODO check type and point, 根据性能等级判定当前point是否需要获取时间
	// TODO write to channel
}

func (pc *perfDist) Stop() {

}

func (pc *perfDist) analysis() {
	// TODO read from channel
	// lock map
}

func (pc *perfDist) perfDump() {

}
