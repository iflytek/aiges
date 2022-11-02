//go:build !linux
// +build !linux

package widget

func NewWidget(plugin string) WidgetInner {
	switch plugin {
	case "python":
		return &WidgetPython{}
	default:
		warn()
		usage()
	}
	return nil
}
