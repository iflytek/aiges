package xsf

import "io"

var monitorInst monitor

type monitor interface {
	Query(map[string]string, io.Writer)
}

func StoreMonitor(in monitor) {
	monitorInst = in
}
