package dotnetdiag

import (
	"fmt"
	"net"
)

// Client implement Diagnostic IPC Protocol client.
// https://github.com/dotnet/diagnostics/blob/main/documentation/design-docs/ipc-protocol.md
type Client struct {
	addr string
	dial Dialer
}

// Dialer establishes connection to the given address. Due to the potential for
// an optional continuation in the Diagnostics IPC Protocol, each successful
// connection between the runtime and a Diagnostic Port is only usable once.
//
// Note that the dialer is OS-specific, refer to documentation for details:
// https://github.com/dotnet/diagnostics/blob/main/documentation/design-docs/ipc-protocol.md#transport
type Dialer func(addr string) (net.Conn, error)

// Option overrides default Client parameters.
type Option func(*Client)

// WithDialer overrides default dialer function with d.
func WithDialer(d Dialer) Option {
	return func(c *Client) {
		c.dial = d
	}
}

// Session represents EventPipe stream of NetTrace data created with
// `CollectTracing` command.
//
// A session is expected to be closed with `StopTracing` call (or `Close`),
// as there is a "run down" at the end of a stream session that transmits
// additional metadata. If the stream is stopped prematurely due to a client
// or server error, the NetTrace stream will be incomplete and should
// be considered corrupted.
type Session struct {
	c    *Client
	conn net.Conn
	ID   uint64
}

// CollectTracingConfig contains supported parameters for CollectTracing command.
type CollectTracingConfig struct {
	// CircularBufferSizeMB specifies the size of the circular buffer used for
	// buffering event data while streaming
	CircularBufferSizeMB uint32
	// Providers member lists providers to turn on for a streaming session.
	// See ETW documentation for a more detailed explanation of Keywords, Filters, and Log Level:
	// https://docs.microsoft.com/en-us/message-analyzer/system-etw-provider-event-keyword-level-settings
	Providers []ProviderConfig
}

// NewClient creates a new Diagnostic IPC Protocol client for the transport
// specified - on Unix/Linux based platforms, a Unix Domain Socket will be used, and
// on Windows, a Named Pipe will be used:
//  - /tmp/dotnet-diagnostic-{%d:PID}-{%llu:disambiguation key}-socket (Linux/MacOS)
//  - \\.\pipe\dotnet-diagnostic-{%d:PID} (Windows)
//
// Refer to documentation for details:
// https://github.com/dotnet/diagnostics/blob/main/documentation/design-docs/ipc-protocol.md#transport
func NewClient(addr string, options ...Option) *Client {
	c := &Client{addr: addr}
	for _, option := range options {
		option(c)
	}
	if c.dial == nil {
		c.dial = DefaultDialer()
	}
	return c
}

// CollectTracing creates a new EventPipe session stream of NetTrace data.
func (c *Client) CollectTracing(config CollectTracingConfig) (s *Session, err error) {
	// Every session has its own IPC connection which cannot be reused for any
	// other purposes; in order to close the connection another connection
	// to be opened - see `StopTracing`.
	conn, err := c.dial(c.addr)
	if err != nil {
		return nil, err
	}
	defer func() {
		// The connection should not be disposed if a session has been created.
		if err != nil {
			_ = conn.Close()
		}
	}()

	p := CollectTracingPayload{
		CircularBufferSizeMB: config.CircularBufferSizeMB,
		Format:               FormatNetTrace,
		Providers:            config.Providers,
	}

	if err = writeMessage(conn, CommandSetEventPipe, EventPipeCollectTracing, p.Bytes()); err != nil {
		return nil, err
	}
	var resp CollectTracingResponse
	if err = readResponse(conn, &resp); err != nil {
		return nil, err
	}

	s = &Session{
		c:    c,
		conn: conn,
		ID:   resp.SessionID,
	}

	return s, nil
}

// StopTracing stops the given streaming session started with CollectTracing.
func (c *Client) StopTracing(sessionID uint64) error {
	conn, err := c.dial(c.addr)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	p := StopTracingPayload{SessionID: sessionID}
	if err := writeMessage(conn, CommandSetEventPipe, EventPipeStopTracing, p.Bytes()); err != nil {
		return err
	}
	var resp StopTracingResponse
	if err := readResponse(conn, &resp); err != nil {
		return err
	}
	if resp.SessionID != sessionID {
		return fmt.Errorf("%w: %x", ErrSessionIDMismatch, resp.SessionID)
	}
	return nil
}

func (s *Session) Read(b []byte) (int, error) {
	return s.conn.Read(b)
}

func (s *Session) Close() error {
	return s.c.StopTracing(s.ID)
}
