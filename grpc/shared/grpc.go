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

func (m *GRPCClient) WrapperCreate(userTag string, params map[string]string) (*proto.Handle, error) {
	return m.client.WrapperCreate(context.Background(), &proto.CreateRequest{
		Sid:    params["sid"],
		Tag:    userTag,
		Params: params,
	})
}

func (m *GRPCClient) WrapperDestroy(handle string) (*proto.Ret, error) {
	return m.client.WrapperDestroy(context.Background(), &proto.Handle{
		Handle: handle,
	})
}

func (m *GRPCClient) WrapperWrite(handle string, usrTag string, params map[string]string, reqData []*proto.RequestData) (*proto.Ret, error) {
	return m.client.WrapperWrite(context.Background(), &proto.WriteMessage{

		Handle: handle,
		Req: &proto.Request{
			Params: params,
			List:   reqData,
			Tag:    usrTag,
		},
	})
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl PyWrapper
}

func (m *GRPCServer) WrapperDestroy(ctx context.Context, handle *proto.Handle) (*proto.Ret, error) {
	//TODO implement me
	panic("implement me")
}

func (m *GRPCServer) WrapperCreate(ctx context.Context, request *proto.CreateRequest) (*proto.Handle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *GRPCServer) WrapperWrite(ctx context.Context, message *proto.WriteMessage) (*proto.Ret, error) {
	//TODO implement me
	panic("implement me")
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
