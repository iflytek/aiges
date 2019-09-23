package ratelimit

import (
	"testing"
	"io/ioutil"
	"fmt"
)

func TestLoadConfigToManager(t *testing.T) {
	f,err:=ioutil.ReadFile("cfg.json")
	if err !=nil{
		panic(err)
	}
	r,err:=LoadRateConfig(f)
	fmt.Println(r[0].AppipLimit().GetKey())
	fmt.Println(r[0].IpLimit().GetKey())
}
