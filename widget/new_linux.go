//go:build linux
// +build linux

package widget

func NewWidget(plugin string) WidgetInner {
	switch plugin {
	case "c":
		return &WidgetC{}
	case "python":
		return &WidgetPython{}
	default:
		usage()
	}
	return nil
}
