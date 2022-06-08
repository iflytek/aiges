package dotnetdiag

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

func DefaultDialer() Dialer {
	return func(addr string) (net.Conn, error) {
		return winio.DialPipe(addr, nil)
	}
}

// DefaultServerAddress returns Diagnostic Server named pipe name for the process given.
// https://github.com/dotnet/diagnostics/blob/main/documentation/design-docs/ipc-protocol.md#transport
func DefaultServerAddress(pid int) string {
	return fmt.Sprintf(`\\.\pipe\dotnet-diagnostic-%d`, pid)
}
