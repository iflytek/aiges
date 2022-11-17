//go:build darwin || windows
// +build darwin windows

package widget

import "github.com/xfyun/aiges/utils"

func NewWidget(plugin string, ch *utils.Coordinator) WidgetInner {
	switch plugin {
	case "python":
		return &WidgetPython{
			ch: ch,
		}
	default:
		warn()
		usage()
	}
	return nil
}
