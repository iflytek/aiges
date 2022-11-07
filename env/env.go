package env

import (
	"log"
	"os"
	"runtime"
	"strconv"
)

var (
	AIGES_ENV_NUAME   int = -1
	AIGES_ENV_VERSION string
	AIGES_PLUGIN_MODE string
	SYSArch           string
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
	if pluginMode := os.Getenv("AIGES_PLUGIN_MODE"); len(pluginMode) > 0 {
		if pluginMode != "c" && pluginMode != "python" {
			log.Fatalln("Not Support This Plugin Mode... should be on of c|python...")
		}
		AIGES_PLUGIN_MODE = pluginMode
	} else {
		// 默认c插件模式
		AIGES_PLUGIN_MODE = "c"

	}
	if goos := os.Getenv("GOOS"); len(goos) > 0 {
		SYSArch = goos
	} else {
		// 默认c插件模式
		SYSArch = runtime.GOOS

	}

}
