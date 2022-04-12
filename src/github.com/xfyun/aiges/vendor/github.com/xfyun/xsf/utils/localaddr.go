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
	var ip_str string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("get ip err")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip_str = ipnet.IP.String()
			}
		}
	}
	return ip_str, nil

}

//return netWorkCardAddress map and err
func GetAddrs() (map[string]string, error) {
	netWorkCardAddrsMap := make(map[string]string)
	netCard, err := net.Interfaces()
	if err != nil {
		return netWorkCardAddrsMap, err
	}
	for _, v := range netCard {
		addrs, err := v.Addrs()
		if err != nil {
			return netWorkCardAddrsMap, err
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					netWorkCardAddrsMap[v.Name] = ipnet.IP.String()
				}
			}
		}
	}
	return netWorkCardAddrsMap, err
}
