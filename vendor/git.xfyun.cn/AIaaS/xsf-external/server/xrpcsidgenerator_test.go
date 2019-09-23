package xsf

import (
	"fmt"
	"testing"
)

func TestXrpcSidGenerator(t *testing.T) {
	var xrpcSid XrpcSidGenerator
	xrpcSid.Init(1, "127.0.0.1", 9165)
	for ix := 0; ix < 10; ix++ {
		sid := xrpcSid.generateSid()
		t.Logf("NO.%d -> sid:%s,len:%d\n", ix, sid, len(sid))
	}
}
func TestParseSid(t *testing.T) {
	var xrpcSid XrpcSidGenerator
	//xrpcSid.ParseSid("00a01cd97c3830000040000105")
	xrpcSid.ParseSid("0ac15d21513e20543ba0025c01")
}
func TestSid2Ip(t *testing.T) {
	var xrpcSid XrpcSidGenerator
	fmt.Println(xrpcSid.Sid2Ip("0ac15d21513e20543ba0025c01"))
}
