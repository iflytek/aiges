package shared

import (
	"errors"
	"fmt"
	"github.com/xfyun/aiges/grpc/proto"
	"golang.org/x/net/context"
	"log"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.WrapperServiceClient }

func (m *GRPCClient) WrapperInit(config map[string]string) error {
	ret, err := m.client.WrapperInit(context.Background(), &proto.InitRequest{
		Config: config,
	})
	if ret.GetRet() != 0 {
		msg := fmt.Sprintf("Call WrapperInit Failed...ret: %d", ret.GetRet())
		log.Println(msg)
		return errors.New(msg)

	}
	return err
}

func (m *GRPCClient) WrapperOnceExec(userTag string, params map[string]string, reqData []*proto.RequestData) (*proto.Response, error) {
	resp, err := m.client.WrapperOnceExec(context.Background(), &proto.Request{
		Params: params,
		List:   reqData,
		Tag:    userTag,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *GRPCClient) Communicate() (proto.WrapperService_CommunicateClient, error) {
	return m.client.Communicate(context.Background())
}

func (m *GRPCClient) WrapperSchema(id string) (*proto.Schema, error) {
	return m.client.WrapperSchema(context.Background(), &proto.SvcId{ServiceId: id})
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl PyWrapper
}

func (m *GRPCServer) WrapperSchema(ctx context.Context, id *proto.SvcId) (*proto.Schema, error) {
	//TODO implement me
	panic("implement me")
}

func (m *GRPCServer) Communicate(server proto.WrapperService_CommunicateServer) error {
	//TODO implement me
	panic("implement me")
}

func (m *GRPCServer) WrapperInit(contenxt context.Context, req *proto.InitRequest) (*proto.Ret, error) {
	err := m.Impl.WrapperInit(req.Config)
	if err != nil {
		return &proto.Ret{Ret: -1}, err

	}
	return &proto.Ret{Ret: 0}, err
}

func (m *GRPCServer) WrapperOnceExec(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	return nil, nil
}

func (m *GRPCServer) TestStream(stream proto.WrapperService_TestStreamServer) error {
	return nil
}
