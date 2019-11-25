package utils

/*
#cgo linux LDFLAGS: -lnuma -lpthread
#include <stdlib.h>
#include <stdio.h>
#include <numa.h>
#include <pthread.h>

// @return error info if setAffinity fail, else NULL;
char* SetProcessAffinityOnNode(int nodeIndex)
{
	int ret = numa_available();
	if (ret < 0){
		return "system doesn't support numa api";
	}

	ret = numa_run_on_node(nodeIndex);
	if (ret < 0){
		return "set process affinity fail";
	}
	return NULL;
}

import "C"
import (
	"errors"
)
*/
func NumaBind(numaIndex int) error {
/*
    if numaIndex >= 0 {
		err := C.SetProcessAffinityOnNode(C.int(numaIndex))
		if err != nil {
			return errors.New("NumaBind fail, " + C.GoString(err))
		}
	}
*/
	return nil
}
