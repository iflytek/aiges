package env

import (
	"os"
	"strconv"
)

var (
	AIGES_ENV_NUAME   int = -1
	AIGES_ENV_VERSION string
)

func Parse() {
	if numa := os.Getenv("AIGES_ENV_NUMA"); len(numa) > 0 {
		if i, _ := strconv.Atoi(numa); i >= 0 {
			AIGES_ENV_NUAME = i
		}
	}

	if ver := os.Getenv("AIGES_WRAPPER_VERSION"); len(ver) > 0 {
		AIGES_ENV_VERSION = ver
	}

}
