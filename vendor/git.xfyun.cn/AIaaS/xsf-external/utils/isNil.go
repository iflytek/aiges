package utils

import "unsafe"

type eFace struct {
	rType unsafe.Pointer
	data  unsafe.Pointer
}

func IsNil(obj interface{}) bool {
	if nil == obj {
		return true
	}
	return (*eFace)(unsafe.Pointer(&obj)).data == nil
}
