package xsf

import (
	"fmt"
	"github.com/xfyun/xsf/utils"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const sidVer int64 = 0

type XrpcSidGenerator struct {
	ver       int64  // 大于等于0，小于10
	ipDec     string // ip的十进制表示
	ipHex     string //ip的十六进制表示
	port      int64  //端口号，小于65535
	count     int64  //递增序列
	startTime int64  // xrpcSidGenerator的初始化时间
}

func NewSidGenerator(ver int64, ipDec string, port int64) (*XrpcSidGenerator, error) {
	if !utils.IsIpv4(ipDec) {
		return nil, fmt.Errorf("ip:%v is not a legal ipv4 addr", ipDec)
	}
	var xrpcSidGenerator XrpcSidGenerator
	xrpcSidGenerator.Init(ver, ipDec, port)
	return &xrpcSidGenerator, nil
}
func (s *XrpcSidGenerator) Init(ver int64, ipDec string, port int64) {
	s.startTime = time.Now().Unix()
	s.ver = ver
	s.ipDec = ipDec
	s.port = port
	ipSplit := strings.Split(s.ipDec, ".")
	ipHex0, _ := strconv.Atoi(ipSplit[0])
	ipHex1, _ := strconv.Atoi(ipSplit[1])
	ipHex2, _ := strconv.Atoi(ipSplit[2])
	ipHex3, _ := strconv.Atoi(ipSplit[3])
	s.ipHex = fmt.Sprintf("%.2x%.2x%.2x%.2x", ipHex0, ipHex1, ipHex2, ipHex3)
	return
}
func (s *XrpcSidGenerator) generateSid() string {
	ver := s.ver
	ipHex := s.ipHex
	port := s.port
	count := atomic.AddInt64(&s.count, 1)
	ts := time.Now().Unix() - s.startTime
	return fmt.Sprintf("%1d%08.8s%04x%06x%07x", ver, ipHex, port, count, ts)
}
func (s *XrpcSidGenerator) ParseSid(sid string) {
	sidBytes := []byte(sid)
	loggerStd.Println(sid, string(sidBytes), )
	ver, ipHex, port, count, ts := sidBytes[:1], sidBytes[1:9], sidBytes[9:13], sidBytes[13:19], sidBytes[19:26]
	loggerStd.Printf("ver:%s, ipHex:%s, port:%s, count:%s, ts:%s\n", ver, ipHex, port, count, ts)
	loggerStd.Printf("ver:%s, ipDec:%s, port:%s, count:%s, ts:%s\n", ver, s.IpHex2Dec(ipHex), s.Hex2Dec(sidBytes[9:13]), s.Hex2Dec(sidBytes[13:19]), s.Hex2Dec(sidBytes[19:26]))
}

func (s *XrpcSidGenerator) Sid2Ip(sid string) string {
	sidBytes := []byte(sid)
	return fmt.Sprintf("%v:%v", s.IpHex2Dec(sidBytes[1:9]), s.Hex2Dec(sidBytes[9:13]))
}
func (s *XrpcSidGenerator) IpHex2Dec(hexData []byte) string {
	if len(hexData) < 8 {
		return ""
	}
	ip1, _ := strconv.ParseInt(string(hexData[0:2]), 16, 64)
	ip2, _ := strconv.ParseInt(string(hexData[2:4]), 16, 64)
	ip3, _ := strconv.ParseInt(string(hexData[4:6]), 16, 64)
	ip4, _ := strconv.ParseInt(string(hexData[6:8]), 16, 64)

	return fmt.Sprintf("%v.%v.%v.%v",
		strconv.Itoa(int(ip1)), strconv.Itoa(int(ip2)), strconv.Itoa(int(ip3)), strconv.Itoa(int(ip4)))
}
func (s *XrpcSidGenerator) Hex2Dec(hexData []byte) string {
	decData, _ := strconv.ParseInt(string(hexData), 16, 64)
	return strconv.Itoa(int(decData))
}
