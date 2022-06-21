package util

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var (
	index        int64  = 0
	Location     string = "hu"
	LocalIP      string = "127.0.0.1"
	ShortLocalIP string = "0001"
	Port         string = "5090"
)

//func init() {
//	n := xsf.Net{}
//	tmp, err := n.GetHostByName("", "")
//	if err != nil || len(tmp) < 1 {
//		panic(err)
//	}
//	LocalIP = tmp[0]
//
//	ip := net.ParseIP(LocalIP)
//	var ipSec3, ipSec4 int
//	if ip != nil {
//		ipSec3 = (int)(ip[14])
//		ipSec4 = (int)(ip[15])
//		ip3 := ipSec3 & 0xFF
//		ip4 := ipSec4 & 0xFF
//		ShortLocalIP = fmt.Sprintf("%02x%02x", ip3, ip4)
//	} else {
//		ShortLocalIP = "xxxx"
//	}
//
//	Port = "5090"
//}

func NewSid(sub string) string {
	if len(sub) == 0 {
		sub = "ase"
	}
	pid := os.Getpid() & 0xFF
	index_now := atomic.AddInt64(&index, 1) & 0xffff
	tmint := time.Now().UnixNano() / 1000000
	tm := fmt.Sprintf("%011x", tmint)
	sid := fmt.Sprintf("%3s%04x%04x@%2s%s%04s%02s0", sub, pid, index_now, Location, tm[len(tm)-11:], ShortLocalIP, Port[:2])
	return sid
}
