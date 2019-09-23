package utils

import "testing"

func TestGetAddr(t *testing.T) {
	if addr, err := GetAddr(); err != nil {
		t.Failed()
	} else {
		t.Logf("addr is : %v\n", addr)
	}
}

func TestGetAddrs(t *testing.T) {
	if netWorkCardAddrsMap, err := GetAddrs(); err != nil {
		t.Failed()
	} else {
		for k, v := range netWorkCardAddrsMap {
			t.Logf("netWorkCard name:%v,addr:%v\n", k, v)
		}
	}
}
