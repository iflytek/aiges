//go:build linux && cgo
// +build linux,cgo

package widget

func NewWidget(plugin string, ch *utils.Coordinator) WidgetInner {
	switch plugin {
	case "c":
		return &WidgetC{}
	case "python":
		return &WidgetPython{
			ch: ch,
		}
	default:
		usage()
	}
	return nil
}
