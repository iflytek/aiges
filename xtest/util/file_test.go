package util

import (
	"testing"
)

func TestCompFunc(t *testing.T) {
	flag := 0 // 不排序
	s1 := "11"
	s2 := "2"
	ans := CompFunc(flag, s1, s2) // return s1 == s2
	if ans {
		t.Fatal()
	}
}

func TestReadDir(t *testing.T) {
	dir := "./"
	ansBytes, err := ReadDir(dir, "e", 1) // 文件升序读取
	if err != nil {
		t.Fatal(err)
	}
	ans := []string{}
	for _, v := range ansBytes {
		ans = append(ans, string(v))
	}
}
