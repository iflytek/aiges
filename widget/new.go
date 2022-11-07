//go:build darwin || windows
// +build darwin windows

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
