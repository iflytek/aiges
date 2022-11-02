// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/xfyun/aiges/grpc/proto"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"wrapper_grpc": &WrapperGRPCPlugin{},
	//"kv":           &KVPlugin{},
}

// KV is the interface that we're exposing as a plugin.
type PyWrapper interface {
	WrapperInit(config map[string]string) error
	WrapperOnceExec(params map[string]string, reqData []*proto.RequestData) (*proto.Response, error)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type WrapperGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl PyWrapper
}

func (p *WrapperGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterWrapperServiceServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *WrapperGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewWrapperServiceClient(c)}, nil
}
