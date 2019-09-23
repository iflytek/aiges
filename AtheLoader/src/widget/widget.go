package widget

import (
	"flag"
	"fmt"
	"os"
	"service"
)

const (
	wrapperC  = "libwrapper.so" // Plugin C
	wrapperGo = "libwrapper.so" // Plugin Go
	wrapperPy = "wrapper.py"    // Plugin Python
)

var pluginMode = flag.String("plugin", "c", "plugin mode, c/go/py is supported")

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
	case "go":
		return &WidgetGo{}
	case "py":
		return &WidgetPy{}
	default:
		usage()
	}
	return nil
}

func usage() {
	fmt.Printf("usage: flag -plugin=c/go/py\n" +
		"-plugin=c		load libwrapper.so (defalut mode)\n" +
		"-plugin=go		load libwrapper.so\n" +
		"-plugin=py		load wrapper.py\n",
	)
	os.Exit(0)
}
