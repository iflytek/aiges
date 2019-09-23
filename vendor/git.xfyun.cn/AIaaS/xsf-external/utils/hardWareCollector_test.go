package utils

import (
	"fmt"
	"testing"
)

func Test_HardWareCollector(t *testing.T) {

	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	inst, err := NewHardWareCollector()
	checkErr(err)

	{
		models, modelsErr := inst.GetModel()
		fmt.Printf("models:%v,modelsErr:%v\n", models, modelsErr)
	}
	{
		cpus, cpusErr := inst.GetCpus()
		fmt.Printf("cpus:%v,cpusErr:%v\n", cpus, cpusErr)
	}
	{
		cpuMHz, cpuMHzErr := inst.GetCpuMHz()
		fmt.Printf("cpuMHz:%v,cpuMHzErr:%v\n", cpuMHz, cpuMHzErr)
	}
	{
		inDocker, inDockerErr := inst.InDocker()
		fmt.Printf("InDocker:%v,inDockerErr:%v\n", inDocker, inDockerErr)
	}
	{
		cpuLimit, cpuLimitErr := inst.GetCpuLimit()
		fmt.Printf("cpuLimit:%v,cpuLimitErr:%v\n", cpuLimit, cpuLimitErr)
	}
	{
		cpuUsage, cpuUsageErr := inst.GetCpuUsage()
		fmt.Printf("cpuUsage:%v,cpuUsageErr:%v\n", cpuUsage, cpuUsageErr)
	}
	{
		memUsage, memUsageErr := inst.GetMemUsage()
		fmt.Printf("memUsage:%v,memUsageErr:%v\n", memUsage, memUsageErr)
	}

}