package util

import (
	"github.com/pterm/pterm"
	"go.uber.org/atomic"
)

// ProgressShow jbzhou5 使用pterm绘制一些进度可视化
func ProgressShow(cnt *atomic.Int64, cnt1 int64) {
	// Create progressbar as fork from the default progressbar.
	p, _ := pterm.DefaultProgressbar.WithTotal(int(cnt1)).WithTitle("Xtest testing ").WithShowCount(true).Start()
	for i, pre := int64(p.Total), int64(p.Total); i > 0; {
		//pterm.Success.Println("Xtest testing " + strconv.Itoa(i)) // If a progressbar is running, each print will be printed above the progressbar.
		pre, i = i, cnt.Load()
		p.Add(int(pre - cnt.Load())) // Increment the progressbar by one. Use Add(x int) to increment by a custom amount.
		p.Add(p.Total - max(p.Current, p.Total))
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
