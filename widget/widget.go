package widget

import (
	"flag"
	"fmt"
	"github.com/xfyun/aiges/service"
	"os"
)

const (
	wrapperC         = "libwrapper.so" // Plugin C
	wrapperGo        = "libwrapper.so" // Plugin Go
	wrapperPythonCmd = "/Users/yangyanbo/anaconda3/envs/aiges-python/bin/python grpc/examples/wrapper-python/plugin.py"
)

var pluginMode = flag.String("plugin", "c", "plugin mode, c/go is supported")

type WidgetInner interface {
	Open() (err error)
	Close()
	Register(srv *service.EngService) (err error)
	Version() (version string)
}

func usage() {
	fmt.Printf("\n")
	os.Exit(0)
}

func warn() {
	fmt.Println("Non linux platform only support python plugin..\n please set env variable by using command " +
		"'export AIGES_PLUGIN_MODE=python/c' ")

}
