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
	fmt.Printf("usage: flag -plugin=c/go/py\n" +
		"-plugin=c		load libwrapper.so (defalut mode)\n" +
		"-plugin=go		load libwrapper.so\n",
	)
	os.Exit(0)
}

func warn() {
	fmt.Println("Non linux platform only support python plugin.. ")

}
