package utils

import (
	"testing"
)

func Test_LookupHost1(t *testing.T) {
	addrs, err := LookupHost("baidu.com", "")
	if err != nil {
		t.Errorf("LookupHost(\"baidu.com\",\"\") failed:%v", addrs)
	}
	LookupHost("baidu.com", "114.114.114.114:53")
	if err != nil {
		t.Errorf("LookupHost(\"baidu.com\",\"114.114.114.114:53\") failed:%v", addrs)
	}
	LookupHost("baidu.com", "1.1.1.1")
	if err == nil {
		t.Errorf("LookupHost(\"baidu.com\",\"0.0.0.0\") failed:%v", addrs)
	}
}
