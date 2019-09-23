/*
## features ##
- support get the special network card address
- support get the first not loopback address
- support get the first not loopback address and ignore some network cards
*/

/*
* @file	xsflog.go
* @brief	get network card address
* @author	sqjian
* @version	1.0
* @date		2017.11.14
*/
package utils

import (
	"net"
	"fmt"
)

//get the first not loopback address
func GetAddr() (string, error) {
	var ipStr string
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		return "", fmt.Errorf("get ip err")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				ipStr = ipnet.IP.String()
			}
		}
	}
	return ipStr, nil

}

//return netWorkCardAddress map and err
func GetAddrs() (map[string]string, error) {
	netWorkCardAddrsMap := make(map[string]string)
	netCard, err := net.Interfaces()
	if nil != err {
		return netWorkCardAddrsMap, err
	}
	for _, v := range netCard {
		addrs, err := v.Addrs()
		if nil != err {
			return netWorkCardAddrsMap, err
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if nil != ipnet.IP.To4() {
					netWorkCardAddrsMap[v.Name] = ipnet.IP.String()
				}
			}
		}
	}
	return netWorkCardAddrsMap, err
}
