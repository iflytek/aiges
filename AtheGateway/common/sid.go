package common

import (
	"os"
	"sync/atomic"
	"time"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"net"
	"strconv"
)



var index int64 = 0

var (
	Location 		string ="dx"
	LocalIP 	string
	ShortLocalIP string
	Port 	string
	PortInt  int
)

func NewSid(sub string) (string, error) {
	if len(sub) == 0 {
		sub = "src"
	}
	pid := os.Getpid() & 0xFF
	index_now := atomic.AddInt64(&index, 1) & 0xffff
	tmint := time.Now().UnixNano()/1000000
	tm := fmt.Sprintf("%011x",tmint)
	sid := fmt.Sprintf("%3s%04x%04x@%2s%s%04s%02s2", sub, pid, index_now, Location,tm[len(tm)-11:] ,ShortLocalIP, Port[:2])
	return sid, nil
}


func InitSidGenerator(inip string,port,location string){
	Location = location
	if inip == ""{
		n := xsf.Net{}
		tmp,err:= n.GetHostByName("","")
		if err != nil || len(tmp) < 1 {
			panic(err)
		}
		LocalIP = tmp[0]
	}else{
		LocalIP = inip
	}
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

	Port = port
	PortInt ,_ = strconv.Atoi(Port)
}

var hexA = []byte{'0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f'}
//  12       1100

func getHexlen(i int) int {
	c:=0
	if i==0{
		return 1
	}
	for i>0{
		i= i/ 16
		c++
	}
	return c
}

func toHexString(i int)string{
	if i==0{
		return "0"
	}
	bit:=getHexlen(i)
	bf:=make([]byte,bit)

	for i>0{
		n := i%16
		bit --
		if bit==-1{
			bit=0
			break
		}
		bf[bit] = hexA[n]
		i= i/ 16
	}

	return string(bf)
}