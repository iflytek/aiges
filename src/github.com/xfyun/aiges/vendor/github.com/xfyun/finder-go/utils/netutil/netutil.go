package netutil

import (
	"errors"
	"log"
	"net"
	"strings"
)

// GetLocalIP get local ip
func GetLocalIP(url string) (string, error) {
	var host string
	var port string
	var localIP string
	items := strings.Split(url, ":")
	if len(items) == 3 {
		host = strings.Replace(items[1], "/", "", -1)
		port = items[2]
	} else if len(items) == 2 {
		host = strings.Replace(items[0], "/", "", -1)
		port = items[1]
	} else {
		host = url
		port = "80"
	}

	if len(host) == 0 {
		return "", errors.New("GetLocalIP:invalid remote url")
	}
	if len(port) == 0 {
		port = "80"
	}
	ips, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		conn, err := net.Dial("tcp", ip+":"+port)
		if err != nil {
			log.Println("GetLocalIP:", err)
			continue
		}
		localIP = conn.LocalAddr().String()
		log.Println("GetLocalIP:ok")
		err = conn.Close()
		if err != nil {
			log.Println("GetLocalIP:", err)
			break
		}
		break
	}
	if len(localIP) == 0 {
		return "", errors.New("GetLocalIP:failed")
	}

	return localIP, nil
}
