package util

import (
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"net"
	"os"
	"sync/atomic"
	"time"
)

var (
	index        int64  = 0
	Location     string = "dx"
	LocalIP      string
	ShortLocalIP string
	Port         string
)

func init() {
	n := xsf.Net{}
	tmp, err := n.GetHostByName("", "")
	if err != nil || len(tmp) < 1 {
		panic(err)
	}
	LocalIP = tmp[0]

	ip := net.ParseIP(LocalIP)
	var ipSec3, ipSec4 int
	if ip != nil {
		ipSec3 = (int)(ip[14])
		ipSec4 = (int)(ip[15])
		ip3 := ipSec3 & 0xFF
		ip4 := ipSec4 & 0xFF
		ShortLocalIP = fmt.Sprintf("%02x%02x", ip3, ip4)
	} else {
		panic("Bad IP !! " + LocalIP)
	}

	Port = "5090"
}

func NewSid(sub string) string {
	if len(sub) == 0 {
		sub = "tst"
	}
	pid := os.Getpid() & 0xFF
	index_now := atomic.AddInt64(&index, 1) & 0xffff
	tmint := time.Now().UnixNano() / 1000000
	tm := fmt.Sprintf("%011x", tmint)
	sid := fmt.Sprintf("%3s%04x%04x@%2s%s%04s%02s0", sub, pid, index_now, Location, tm[len(tm)-11:], ShortLocalIP, Port[:2])
	return sid
}
