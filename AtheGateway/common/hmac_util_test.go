package common

import (
	"testing"
	"sync"
	"fmt"
)

func TestHmacWithShaTobase64(t *testing.T) {
	type args struct {
		algorithm string
		data      string
		key       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HmacWithShaTobase64(tt.args.algorithm, tt.args.data, tt.args.key); got != tt.want {
				t.Errorf("HmacWithShaTobase64() = %v, want %v", got, tt.want)
			}
		})
	}
}
var total = 10000000
func TestMap(t *testing.T) {
	m:=sync.Map{}
	m.Store("haha","sdfds")
	m.Store("haha2","sdfds")
	m.Store("haha3","sdfds")
	m.Store("haha4","sdfds")
	m.Store("haha5","sdfds")
	m.Store("hadfha5","sdfds")
	m.Store("hahgfga5","sdfds")
	m.Store("hgfgaha5","sdfds")
	m.Store("hagfha5","sdfds")
	m.Store("hahggfa5","sdfds")
	m.Store("hahdda5","sdfds")
	m.Store("hahagf5","sdfds")
	m.Store("hahggda5","sdfds")
var s interface{}
	for i:=0;i<total;i++{
		s,_ = m.Load("haha")
	}
	fmt.Println(s)

}

func TestMap2(t *testing.T) {
	m:=make(map[string]interface{})
	m["haha"]="sdfds"
	m["haha2"]="sdfds"
	m["haha3"]="sdfds"
	m["haha4"]="sdfds"
	m["df"]="sdfds"
	m["hafdfha5"]="sdfds"
	m["hafdfha5"]="sdfds"
	m["hahffa5"]="sdfds"
	m["hahdfa5"]="sdfds"
	m["hahdda5"]="sdfds"
	m["hahssa5"]="sdfds"
	m["haha5"]="sdfds"
	m["hasgha5"]="sdfds"
	m["hahssa5"]="sdfds"
	var s interface{}

	for i:=0;i<total;i++{
		s = m["haha"]
	}

	fmt.Println(s)
}

func TestMap3(t *testing.T) {
	m:=Map{m:map[interface{}]interface{}{}}
	m.Put("haha","sdfds")
	m.Put("haha2","sdfds")
	m.Put("haha3","sdfds")
	m.Put("haha4","sdfds")
	m.Put("haha5","sdfds")
	var s interface{}

	for i:=0;i<total;i++{
		s = m.Get("haha")
	}

	
	fmt.Println(s)
}

type Map struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func (m *Map)Put(k,v interface{})  {
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}

func (m *Map)Get(key interface{})interface{}  {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}