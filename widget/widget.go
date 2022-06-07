package widget

import (
	"flag"
	"fmt"
	"github.com/xfyun/aiges/service"
	"os"
)

const (
	wrapperC  = "libwrapper.so" // Plugin C
	wrapperGo = "libwrapper.so" // Plugin Go
)

var pluginMode = flag.String("plugin", "c", "plugin mode, c/go is supported")

type WidgetInner interface {
	Open() (err error)
	Close()
	Register(srv *service.EngService) (err error)
	Version() (version string)
}

func NewWidget() WidgetInner {
	switch *pluginMode {
	case "c":
		return &WidgetC{}
	default:
		usage()
	}
	return nil
}

func usage() {
	fmt.Printf("usage: flag -plugin=c/go/py\n" +
		"-plugin=c		load libwrapper.so (defalut mode)\n" +
		"-plugin=go		load libwrapper.so\n",
	)
	os.Exit(0)
}
