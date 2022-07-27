package util

import (
	"testing"
	"time"
)

func TestCharts_Draw(t *testing.T) {
	c := Charts{
		Vals: LinesData{
			Title: "Resource Record",
			BarValues: []LineYValue{{"asd", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"hgj", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"dfg45r", []float64{1, 2, 700, 100, 200, 6, 700}},
				{"2342sr", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"das21-asd", []float64{300000, 200000, 400000, 100000, 400000, 450000, 400000}},
				{"csc", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"mhj", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"876ijgh", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"fbfdv", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"67ds", []float64{400, 10000, 200, 50, 5, 800, 7}},
				{"67bdfv", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"sdf324", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"vdf67", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"vdfs234", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"123sdf", []float64{1, 2, 700, 100, 200, 6, 700}},
				{"aasdasd", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"aasd", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"basd", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"cczx", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"qweqw", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"asdadf", []float64{1, 2, 300, 100, 200, 6, 700}},
				{"fghfh", []float64{400, 500000, 200, 50, 5, 800, 7}},
				{"erttyrt", []float64{1, 2, 300, 100, 200, 6, 700}},
			},
		},
		Dst:     "lines.png",
		XValues: nil,
	}

	for i := 0; i < 7; i++ {
		c.XValues = append(c.XValues, float64(time.Now().Add(time.Second).Unix())+float64(i))
	}

	err := c.Draw()
	if err != nil {
		t.Fatal(err)
	}
}
