package daemon

/*
* @file	mssSidGenerator.go
* @brief	mss sid generator
* @author	sqjian
* @version	1.0
* @date		2017.11.27
*/
import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type MssSidGenerator struct {
	sub   string
	ipDec string
	ipHex string
	area  string
	count int64
}

func (s *MssSidGenerator) Init(sub string, ip string, area string) {
	s.ipDec = ip
	s.sub = sub
	s.area = area
	ipSplit := strings.Split(ip, ".")
	ipHex0, _ := strconv.Atoi(ipSplit[0])
	ipHex1, _ := strconv.Atoi(ipSplit[1])
	ipHex2, _ := strconv.Atoi(ipSplit[2])
	ipHex3, _ := strconv.Atoi(ipSplit[3])
	s.ipHex = fmt.Sprintf("%.2x%.2x%.2x%.2x", ipHex0, ipHex1, ipHex2, ipHex3)
	return
}
func (s *MssSidGenerator) GenerateSid(prefix string) (sid string) {
	fill8 := fmt.Sprintf("%08.8s", fmt.Sprintf("%08x", atomic.AddInt64(&s.count, 1)))
	fill4 := fmt.Sprintf("%04.04x", fmt.Sprintf("%04x", os.Getpid()))
	baseTime, _ := time.Parse("2006-01-02 15:04:05.000", "2010-10-01 00:00:00.000")
	timeNow, _ := time.Parse("2006-01-02 15:04:05.000", time.Now().Format("2006-01-02 15:04:05.000"))
	timeNum := timeNow.Sub(baseTime).Seconds()
	timeInt := int(timeNum)
	timeStamp := fmt.Sprintf("%08.8s", fmt.Sprintf("%08x", timeInt))
	return prefix + "@" + s.sub + fill8 + "@" + s.area + fill4 + timeStamp + s.ipHex + "00"
}
