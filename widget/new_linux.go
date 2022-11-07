//go:build linux && cgo
// +build linux,cgo

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
