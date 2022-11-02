package shared

import (
	"github.com/xfyun/aiges/grpc/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.WrapperServiceClient }

func (m *GRPCClient) WrapperInit(config map[string]string) error {
	_, err := m.client.WrapperInit(context.Background(), &proto.InitRequest{
		Config: config,
	})
	return err
}

func (m *GRPCClient) WrapperOnceExec(params map[string]string, reqData []*proto.RequestData) (*proto.Response, error) {
	resp, err := m.client.WrapperOnceExec(context.Background(), &proto.OnceExecRequest{
		Params:  params,
		ReqData: reqData,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl PyWrapper
}

func (m *GRPCServer) WrapperInit(contenxt context.Context, req *proto.InitRequest) (*proto.Ret, error) {
	err := m.Impl.WrapperInit(req.Config)
	if err != nil {
		return &proto.Ret{Ret: -1}, err

	}
	return &proto.Ret{Ret: 0}, err
}

func (m *GRPCServer) WrapperOnceExec(ctx context.Context, req *proto.OnceExecRequest) (*proto.Response, error) {
	return nil, nil
}

func (m *GRPCServer) TestStream(stream proto.WrapperService_TestStreamServer) error {
	return nil
}
