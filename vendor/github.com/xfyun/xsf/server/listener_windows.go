package xsf

import (
	"fmt"
	"net"
	"runtime"
)

const (
	CONNTYPE = "tcp"
)

func NewListener(report int, addr string) (l net.Listener, err error) {
	switch goos := runtime.GOOS; goos {
	case "windows":
		{
			return net.Listen(CONNTYPE, addr)
		}
	default:
		{
			return nil, fmt.Errorf("This operating system %v is not supported.\n", goos)
		}
	}
	return net.Listen("tcp", addr)
}
