package xsf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xfyun/xsf/utils"
	"google.golang.org/grpc"
	"io"
	"math"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

/*
/debug/pprof
/debug/pprof/profile
/debug/pprof/heap
/debug/pprof/block
/debug/pprof/goroutine
/debug/pprof/threadcreate
/debug/pprof/mutex

http://127.0.0.1:8080/debug/pprof/block
http://127.0.0.1:8080/debug/pprof/profile
http://127.0.0.1:8080/debug/pprof/heap
http://127.0.0.1:8080/debug/pprof/goroutine
http://127.0.0.1:8080/debug/pprof/threadcreate
http://127.0.0.1:8080/debug/pprof/mutex

pprof -http "localhost:80"  http://127.0.0.1:8080/debug/pprof/block
*/

func healthServer(w io.Writer) {
	rst := make(map[string]error)
	HealthMu.RLock()
	for k, v := range AdminCheckList {
		rst[k] = v.Check()
	}
	HealthMu.RUnlock()
	resp := func(rst map[string]error) string {
		var res []string
		res = append(res, fmt.Sprintf("health staus:"))
		for k, v := range rst {
			res = append(res, fmt.Sprintf("%v\tcheck: %v", k, v))
		}
		return strings.Join(res, "\n")
	}(rst)
	io.WriteString(w, resp)
}
func statusServer(w io.Writer) {
	var res []string
	res = append(res, fmt.Sprintf("pid:\t\t%v", p.GetPid()))
	user, _ := p.GetUser()
	res = append(res, fmt.Sprintf("user:\t\t%v", user))
	res = append(res, fmt.Sprintf("uptime:\t\t%vs", p.GetUptime()))
	res = append(res, fmt.Sprintf("goroutineid:\t%v", p.GetGoroutineID()))
	res = append(res, fmt.Sprintf("goroutines:\t%v", p.GetGoroutines()))
	resp := strings.Join(res, "\n")
	io.WriteString(w, resp)
}
func goroutineServer(w io.Writer) {
	ProcessInput("lookup goroutine", w)
}
func heapServer(w io.Writer) {
	ProcessInput("lookup heap", w)
}
func threadcreateServer(w io.Writer) {
	ProcessInput("lookup threadcreate", w)
}
func blockServer(w io.Writer) {
	ProcessInput("lookup block", w)
}
func gcsummaryServer(w io.Writer) {
	ProcessInput("gc summary", w)
}

// 生成 CPU 报告
func cpuProfile(w io.Writer) {
	if err := pprof.StartCPUProfile(w); err != nil {
		return
	}
	defer pprof.StopCPUProfile()

	time.Sleep(30 * time.Second) //硬编码，后续优化
}
func filterNumber(in int64) int64 {
	if in == math.MaxInt64 {
		return 0
	}
	if in == math.MinInt64 {
		return 0
	}
	return in
}

// 生成 metrics 报告
func metrics(rawQuery map[string]string, w io.Writer) {
	var timeUnit time.Duration

	{
		timeUnitStr, timeUnitStrOk := rawQuery[MetricsTimeUnit]
		if timeUnitStrOk {
			timeUnitInt, timeUnitIntErr := strconv.Atoi(timeUnitStr)
			if timeUnitIntErr == nil {
				timeUnit = time.Millisecond * time.Duration(timeUnitInt)
			} else {
				timeUnit = defaultMetricsTimeUnit
			}
		} else {
			timeUnit = defaultMetricsTimeUnit
		}
	}

	{
		//被动触发统计
		for _, slidingErrCodeWindowsItem := range slidingErrCodeWindows {
			errCodeMap := slidingErrCodeWindowsItem.getStats(timeUnit)
			for errCode, errCodeCnt := range errCodeMap {
				registryInst.WithLabelValues(
					float64(errCodeCnt),
					WithMetricsClass(REQUEST),
					WithMetricsType("passive"),
					WithMetricsCode(strconv.Itoa(int(errCode))))
			}
		}

		for _, slidingDelayWindowsItem := range slidingDelayWindows {
			max, min, avg, qps := slidingDelayWindowsItem.getStats()
			registryInst.WithLabelValues(
				float64(filterNumber(max)),
				WithMetricsClass(DELAY),
				WithMetricsType("max"))
			registryInst.WithLabelValues(
				float64(filterNumber(min)),
				WithMetricsClass(DELAY),
				WithMetricsType("min"))
			registryInst.WithLabelValues(
				float64(filterNumber(avg)),
				WithMetricsClass(DELAY),
				WithMetricsType("avg"))
			registryInst.WithLabelValues(
				float64(filterNumber(qps)),
				WithMetricsClass(DELAY),
				WithMetricsType("qps"))
		}
		registryInst.WithLabelValues(
			float64(atomic.LoadInt64(&ConcurrentStatistics)),
			WithMetricsClass(CONCURRENT))
	}
	_, _ = w.Write(registryInst.toBytes())
}

// 生成堆内存报告
func heapProfile(w io.Writer) {
	_ = pprof.WriteHeapProfile(w)
}

// 生成追踪报告
func traceProfile(w io.Writer) {
	if err := trace.Start(w); err != nil {
		return
	}
	defer trace.Stop()

	time.Sleep(30 * time.Second) //硬编码，后续优化
}

//用于确认服务启动ok
func healthCheck(target string) bool {
	conn, connErr := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Millisecond*100),
		grpc.WithBlock())
	if connErr != nil {
		loggerStd.Printf("did not connect: %v\n", connErr)
		return false
	}
	defer conn.Close()

	c := utils.NewToolBoxClient(conn)

	query, _ := json.Marshal(map[string]string{"cmd": "health"})
	header, _ := json.Marshal(map[string]string{"method": "GET"})

	cmdServerResp, cmdServerRespErr := c.Cmdserver(context.Background(), &utils.Request{Query: string(query), Headers: string(header), Body: ""})
	if cmdServerRespErr != nil {
		return false
	}
	if strings.Contains(strings.ToLower(string(cmdServerResp.Body)), "err") {
		return false
	}
	return true
}
func cmdServerRouter(cmd string, rawQuery map[string]string, w io.Writer) {
	switch cmd {
	case "status":
		{
			statusServer(w)
		}
	case "health":
		{
			if reportFlag == reportFailure {
				_, _ = io.WriteString(w, "report error")
				break
			}
			healthServer(w)
		}
	case "goroutine":
		{
			goroutineServer(w)
		}
	case "heap":
		{
			heapServer(w)
		}
	case "threadcreate":
		{
			threadcreateServer(w)
		}
	case "block":
		{
			blockServer(w)
		}
	case "gcsummary":
		{
			gcsummaryServer(w)
		}
	case "rawHeap":
		{
			heapProfile(w)
		}
	case "rawTrace":
		{
			traceProfile(w)
		}
	case "rawCpu":
		{
			cpuProfile(w)
		}
	case "metrics":
		{
			metrics(rawQuery, w)
		}

	case "query":
		{
			if nil == monitorInst {
				w.Write([]byte("query method not defined"))
			} else {
				monitorInst.Query(rawQuery, w)
			}
		}
	default:
		{
			w.Write([]byte("don't support the cmd"))
		}

	}
}
