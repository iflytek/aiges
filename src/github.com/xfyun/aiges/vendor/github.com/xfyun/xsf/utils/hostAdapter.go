package utils

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

//如果host为32位ipv4地址，直接返回
//如果host为域名，则解析host对应的ip地址，然后返回
func HostAdapter(host string, netcard string) (string, error) {
	if IsIpv4(host) {
		return host, nil
	}
	return Host2Ip(host, netcard)
}
func Host2Ip(host string, netcard string) (string, error) {
	/*
		1、如果host存在，则取host对应的ip
	*/
	if host != "" {
		addrs, err := net.LookupHost(host)
		if len(addrs) == 0 {
			return "", fmt.Errorf("can't convert host -> %v to ip", host)
		}
		return addrs[0], err
	}

	/*
		1、如果host不存在，netcard存在，则取netcard对应的ip
	*/
	if host == "" && netcard != "" {
		netCardIp, netErr := func(netcard string) (string, error) {
			ipMap, ipMapErr := GetAddrs()
			return ipMap[netcard], ipMapErr
		}(netcard)
		if netErr != nil {
			return "", netErr
		}
		return netCardIp, nil
	}

	/*
	   如果是 windows，走 dns 取正在用的 ip
	*/
	if runtime.GOOS == "windows" {
		return getWinIp()
	}

	/*
		1、如果host、和netcard都没传，则
			a、去本机hostname
			b、调用LookupHost查找hostname对应的ip
	*/
	if host == "" && netcard == "" {
		hostname, hostnameErr := os.Hostname()
		if hostnameErr != nil {
			return "", hostnameErr
		}
		addrs, err := net.LookupHost(hostname)
		if len(addrs) == 0 {
			return "", fmt.Errorf("can't convert host -> %v to ip", host)
		}
		return addrs[0], err
	}

	addrs, err := net.LookupHost(host)
	if len(addrs) == 0 {
		return "", fmt.Errorf("can't convert host -> %v to ip", host)
	}
	return addrs[0], err
}

//检查该ip是否为合法的ipv4地址
func IsIpv4(ip string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	return true
}
func getWinIp() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	if conn == nil {
		return "", fmt.Errorf("nil conn")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
