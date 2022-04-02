package catch

import (
	"runtime/debug"
)

// go协程recover, for each goroutine; TODO 重载 go func()
func RecoverHandle() {
	if switchOn {
		if err := recover(); err != nil {
			dump(nil, debug.Stack(), err.(error))
		}
	}
}
