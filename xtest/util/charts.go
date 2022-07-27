package util

import (
	"errors"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"math/rand"
	"os"
	"time"
)

const (
	lineChartXAxisName = "Time"
	lineChartYAxisName = "Percentage"
	lineChartHeight    = 700
	lineChartWidth     = 1280
	colorMultiplier    = 256
)

var (
	lineChartStyle = chart.Style{
		Padding: chart.Box{
			Top:  30,
			Left: 150,
		},
	}
	timeFormat = GetHMS
)

type Charts struct {
	Vals    LinesData
	Dst     string    // 保存文件
	XValues []float64 // X轴时间戳
}

type LineYValue struct {
	Name   string
	Values []float64
}

type LinesData struct {
	Title     string
	BarValues []LineYValue
}

// createLineChart 创建线性图
func (c *Charts) createLineChart(title string, xValues []float64, values []LineYValue) error {
	if len(values) == 0 {
		return errors.New("Y axis length is zero")
	}
	// 1、计算X轴
	// X轴内容xValues 及 X轴坐标ticks
	var ticks []chart.Tick
	for _, t := range xValues {
		ticks = append(ticks, chart.Tick{Value: t, Label: timeFormat(t)})
	}
	// 2、生成Series
	var series []chart.Series
	for _, yValue := range values {
		mainSeries := chart.ContinuousSeries{
			Name:    yValue.Name,
			XValues: xValues,
			YValues: yValue.Values,
			Style: chart.Style{StrokeColor: drawing.Color{
				R: uint8(rand.Intn(colorMultiplier)),
				G: uint8(rand.Intn(colorMultiplier)),
				B: uint8(rand.Intn(colorMultiplier)),
				A: uint8(colorMultiplier - 1), // 透明度
			}},
		}
		series = append(series, mainSeries)
	}

	// 3、新建图形
	graph := chart.Chart{
		Title:      title, // 定义图标title名
		Background: lineChartStyle,
		Width:      lineChartWidth,  // 图标宽
		Height:     lineChartHeight, // 图标长
		XAxis: chart.XAxis{
			Name:           lineChartXAxisName, // 定义x轴名称
			ValueFormatter: timeFormat,         // 格式化
		},
		YAxis: chart.YAxis{
			Name: lineChartYAxisName, // 定义y轴名称
		},
		Series: series,
	}
	graph.Elements = []chart.Renderable{chart.LegendLeft(&graph)}
	f, _ := os.Create(c.Dst)
	defer f.Close()
	err := graph.Render(chart.PNG, f)
	return err
}

// Draw 传入绘制数据，绘制条形图
func (c *Charts) Draw() error {
	return c.createLineChart(c.Vals.Title, c.XValues, c.Vals.BarValues)
}

// GetHMS 格式化时间获取时分秒
func GetHMS(v interface{}) string {
	//t, _ := time.ParseInLocation("2006-04-01 11:22:22 000", ts, time.Local)
	//ms := int64(v.(float64)) // millsecond数目
	//t := time.UnixMilli(ms)
	//h, m, s, ns := t.Hour(), t.Minute(), t.Second(), t.Nanosecond()
	//return fmt.Sprintf("%d:%d:%d %d", h, m, s, ns)
	return ""
}

// getNsec 获取纳秒数
func getNsec(cur time.Time) float64 {
	return float64(cur.Unix() * int64(time.Second))
}
