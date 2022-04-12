package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	InterErr = errors.New("inter error")
)

func GetCpuAndMem(topData string) (cpu, mem float64, err error) {
	var data []string
	for _, v := range strings.Split(topData, " ") {
		if strings.Replace(v, " ", "", -1) != "" {
			data = append(data, v)
		}
	}
	mem, err = strconv.ParseFloat(data[len(data)-3], 10)
	if err != nil {
		return cpu, mem, err
	}
	cpu, err = strconv.ParseFloat(data[len(data)-4], 10)
	return cpu, mem, err
}

type cpuCollector struct {
	info []cpu.InfoStat
}

func newCpuCollector() (*cpuCollector, error) {
	info, infoErr := cpu.Info()
	if infoErr != nil {
		return nil, infoErr
	}
	return &cpuCollector{info: info}, nil
}

func (c *cpuCollector) GetModel() ([]string, error) {
	var modes []string
	if len(c.info) == 0 {
		return nil, InterErr
	}
	for _, v := range c.info {
		modes = append(modes, v.ModelName)
	}
	return modes, nil
}

func (c *cpuCollector) GetCpus() (int, error) {
	return cpu.Counts(true)
}

func (c *cpuCollector) GetCpuMHz() ([]float64, error) {
	if len(c.info) == 0 {
		return nil, InterErr
	}
	var rst []float64
	for _, v := range c.info {
		rst = append(rst, v.Mhz)
	}
	return rst, nil
}

func (c *cpuCollector) InDocker() (bool, error) {
	readLines := func(data []byte) ([]string, error) {
		var rst []string
		r := bufio.NewReader(bytes.NewReader(data))
		for {
			c, _, e := r.ReadLine()
			if e == io.EOF {
				return rst, nil
			}
			if e != nil && e != io.EOF {
				return nil, e
			}
			rst = append(rst, string(c))
		}
	}

	markedDocker := func(data []string) (bool, error) {
		for _, dataItem := range data {
			if strings.Contains(dataItem, "name") {
				continue
			}
			if !strings.Contains(dataItem, "docker") {
				return false, errors.New(dataItem)
			}
		}
		return true, nil
	}

	r, e := ioutil.ReadFile("/proc/1/cgroup")
	if e != nil {
		return false, e
	}

	lines, linesErr := readLines(r)
	if linesErr != nil {
		return false, e
	}

	return markedDocker(lines)

}

func (c *cpuCollector) GetCpuLimit() (float64, error) {
	const cfsQuotaUsPath = `/sys/fs/cgroup/cpu/cpu.cfs_quota_us`
	const cfsPeriodUsPath = `/sys/fs/cgroup/cpu/cpu.cfs_period_us`

	readCpuLimit := func(inCfsQuotaUs, inCfsPeriodUs string) (outCfsQuotaUs, outCfsPeriodUs int, err error) {
		cfsQuotaUs, cfsQuotaUsErr := ioutil.ReadFile(inCfsQuotaUs)
		if cfsQuotaUsErr != nil {
			return 0, 0, cfsQuotaUsErr
		}
		outCfsQuotaUs, err = strconv.Atoi(strings.Replace(string(cfsQuotaUs), "\n", "", -1))
		if err != nil {
			return outCfsQuotaUs, 0, err
		}

		cfsPeriodUs, cfsPeriodUsErr := ioutil.ReadFile(inCfsPeriodUs)
		if cfsPeriodUsErr != nil {
			return 0, 0, cfsPeriodUsErr
		}
		outCfsPeriodUs, err = strconv.Atoi(strings.Replace(string(cfsPeriodUs), "\n", "", -1))
		return outCfsQuotaUs, outCfsPeriodUs, err
	}

	outCfsQuotaUs, outCfsPeriodUs, err := readCpuLimit(cfsQuotaUsPath, cfsPeriodUsPath)
	if err != nil {
		return 0, err
	}
	return float64(outCfsQuotaUs) / float64(outCfsPeriodUs), nil
}
func (c *cpuCollector) GetCpuUsage(opt ...interface{}) (cpu float64, err error) {
	var interval time.Duration
	if 0 != len(opt) {
		intervalT, intervalOk := opt[1].(time.Duration)
		if intervalOk {
			interval = intervalT
		}
	}
	if interval > 0 {
		cpu, err = c.GetCpuUsageWithProc(interval)
		dbgLoggerStd.Printf("fn:GetCpuUsage,action:GetCpuUsageWithProc,interval:%v,cpu:%v,err:%v\n",
			interval, cpu, err)
	} else {
		cpu, err = c.GetCpuUsageWithExec()
		dbgLoggerStd.Printf("fn:GetCpuUsage,action:GetCpuUsageWithExec,cpu:%v,err:%v\n",
			cpu, err)
	}
	return cpu, err
}
func (c *cpuCollector) GetCpuUsageWithProc(interval time.Duration) (cpuF float64, err error) {
	if interval == 0 {
		return -1, errors.New("interval is empty")
	}
	getCpuTs := func() (
		cpuTotalSlice float64,
		cpuProcessTotalSlice float64,
	) {
		getFileContent := func(path string) string {
			r, e := ioutil.ReadFile(path)
			if e != nil {
				panic(e)
			}
			return strings.TrimPrefix(
				strings.Replace(
					string(r),
					"  ",
					" ",
					-1,
				),
				" ",
			)
		}
		toFloat := func(in string) float64 {
			v, e := strconv.ParseFloat(in, 64)
			if e != nil {
				panic(e)
			}
			return v
		}
		{
			line := strings.Split(
				getFileContent(
					fmt.Sprintf(`/proc/%v/stat`, os.Getpid()),
				),
				" ",
			)
			cpuProcessTotalSlice = toFloat(line[13]) + toFloat(line[14]) + toFloat(line[15]) + toFloat(line[16])
		}
		{
			r := bufio.NewReader(
				strings.NewReader(getFileContent(`/proc/stat`)),
			)
			line, err := r.ReadString('\n')
			if err != nil {
				panic(err)
			}
			vals := strings.Split(
				strings.Trim(line, "\n"),
				" ",
			)
			for i := 1; i < len(vals); i++ {
				cpuTotalSlice += toFloat(vals[i])
			}
		}
		return
	}
	cpuTotalSlice1, cpuProcessTotalSlice1 := getCpuTs()
	time.Sleep(interval)
	cpuTotalSlice2, cpuProcessTotalSlice2 := getCpuTs()
	return (cpuProcessTotalSlice2 - cpuProcessTotalSlice1) / (cpuTotalSlice2 - cpuTotalSlice1), nil
}

func (c *cpuCollector) GetCpuUsageWithExec() (cpu float64, err error) {
	topRst, topErr := exec.Command("top", "-p 1", "-b", "-n 1").Output()
	if topErr != nil {
		return 0, topErr
	}
	r := bufio.NewReader(strings.NewReader(string(topRst)))
	pidTrue := false
	for {
		c, _, e := r.ReadLine()
		if e == io.EOF {
			break
		}
		if e != nil && e != io.EOF {
			return 0, e
		}
		if strings.Contains(string(c), "PID") {
			pidTrue = true
			continue
		}
		if pidTrue {
			cpu, _, err = GetCpuAndMem(string(c))
			return cpu, err
		}
	}
	return
}

type memCollector struct {
}

func NewMemCollector() (*memCollector, error) {
	return &memCollector{}, nil
}
func (m *memCollector) GetMemUsage() (mem float64, err error) {
	topRst, topErr := exec.Command("top", "-p 1", "-b", "-n 1").Output()
	if topErr != nil {
		return 0, topErr
	}
	r := bufio.NewReader(strings.NewReader(string(topRst)))
	pidTrue := false
	for {
		c, _, e := r.ReadLine()
		if e == io.EOF {
			break
		}
		if e != nil && e != io.EOF {
			return 0, e
		}
		if strings.Contains(string(c), "PID") {
			pidTrue = true
			continue
		}
		if pidTrue {
			_, mem, err = GetCpuAndMem(string(c))
			return mem, err
		}
	}
	return
}

type HardWareCollector struct {
	*cpuCollector
	*memCollector
}

func NewHardWareCollector() (*HardWareCollector, error) {
	c := &HardWareCollector{}
	err := c.Init()
	return c, err
}
func (h *HardWareCollector) Init() (err error) {
	h.cpuCollector, err = newCpuCollector()
	if err != nil {
		return err
	}
	h.memCollector, err = NewMemCollector()

	return err
}
func HardWareCollectorExample() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	inst, err := NewHardWareCollector()
	checkErr(err)

	{
		models, modelsErr := inst.GetModel()
		loggerStd.Printf("models:%v,modelsErr:%v\n", models, modelsErr)
	}
	{
		cpus, cpusErr := inst.GetCpus()
		loggerStd.Printf("cpus:%v,cpusErr:%v\n", cpus, cpusErr)
	}
	{
		cpuMHz, cpuMHzErr := inst.GetCpuMHz()
		loggerStd.Printf("cpuMHz:%v,cpuMHzErr:%v\n", cpuMHz, cpuMHzErr)
	}
	{
		inDocker, inDockerErr := inst.InDocker()
		loggerStd.Printf("InDocker:%v,inDockerErr:%v\n", inDocker, inDockerErr)
	}
	{
		cpuLimit, cpuLimitErr := inst.GetCpuLimit()
		loggerStd.Printf("cpuLimit:%v,cpuLimitErr:%v\n", cpuLimit, cpuLimitErr)
	}
	{
		cpuUsage, cpuUsageErr := inst.GetCpuUsage()
		loggerStd.Printf("cpuUsage:%v,cpuUsageErr:%v\n", cpuUsage, cpuUsageErr)
	}
	{
		memUsage, memUsageErr := inst.GetMemUsage()
		loggerStd.Printf("memUsage:%v,memUsageErr:%v\n", memUsage, memUsageErr)
	}
}
