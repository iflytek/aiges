package analy

import (
	"fmt"
	"regexp"
	"testing"
)

func TestFormatRecord(t *testing.T) {
	record := "id:ase00bb0064@hu179f4fe793a0001500,cost:437,begin:2021-06-10 16:16:28.475002412 +0800 CST m=+0.060459613,end:2021-06-10 16:16:28.475439791 +0800 CST m=+0.060896992\n"
	reg := regexp.MustCompile("id:(.*?),cost:(.*?),begin:(.*?),end:(.*?)\n")
	fmt.Println(reg.FindStringSubmatch(record))
}
