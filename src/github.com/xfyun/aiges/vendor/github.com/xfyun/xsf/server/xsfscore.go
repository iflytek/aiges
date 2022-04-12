package xsf

import (
	"context"
	"github.com/xfyun/xsf/server/internal/bvt"
	"github.com/xfyun/xsf/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"math"
	"net"
	"os"
	"sync"
	"time"
)

type grpcOpt struct {
	maxConcurrentStreams  uint32
	maxReceiveMessageSize int
	maxSendMessageSize    int
	initialWindowSize     int32
	initialConnWindowSize int32
	writeBufferSize       int
	readBufferSize        int
	keepaliveTime         time.Duration
	keepaliveTimeout      time.Duration
}

func (g *grpcOpt) setMaxConcurrentStreams(in uint32) {
	g.maxConcurrentStreams = in
}
func (g *grpcOpt) setMaxReceiveMessageSize(in int) {
	g.maxReceiveMessageSize = in
}
func (g *grpcOpt) setMaxSendMessageSize(in int) {
	g.maxSendMessageSize = in
}
func (g *grpcOpt) setInitialWindowSize(in int32) {
	g.initialWindowSize = in
}
func (g *grpcOpt) setInitialConnWindowSize(in int32) {
	g.initialConnWindowSize = in
}
func (g *grpcOpt) setWriteBufferSize(in int) {
	g.writeBufferSize = in
}
func (g *grpcOpt) setReadBufferSize(in int) {
	g.readBufferSize = in
}
func (g *grpcOpt) setKeepaliveTime(in time.Duration) {
	g.keepaliveTime = in
}
func (g *grpcOpt) setKeepaliveTimeout(in time.Duration) {
	g.keepaliveTimeout = in
}

const maxConcurrentStreams uint32 = 0
const maxReceiveMessageSize = 1024 * 1024 * 4
const maxSendMessageSize = math.MaxInt32
const initialWindowSize int32 = 0
const initialConnWindowSize int32 = 0
const writeBufferSize = 1024 * 1024 * 2
const readBufferSize = 0
const grpcSleepTime = time.Second * 3

var grpcOptInst = &grpcOpt{
	maxConcurrentStreams:  maxConcurrentStreams,
	maxReceiveMessageSize: maxReceiveMessageSize,
	maxSendMessageSize:    maxSendMessageSize,
	initialWindowSize:     initialWindowSize,
	initialConnWindowSize: initialConnWindowSize,
	writeBufferSize:       writeBufferSize,
	readBufferSize:        readBufferSize,
}

type xsfServer struct {
	grpcserver *grpc.Server
}

func (x *xsfServer) Closeout() {
	loggerStd.Println("about to stop grpcserver")
	x.grpcserver.GracefulStop()
}
func (x *xsfServer) run(bc BootConfig, listener net.Listener, srv utils.XsfCallServer) (resErr error) {
	var opts []grpc.ServerOption
	//----------------------------------------------------
	// 注册interceptor
	var interceptor grpc.UnaryServerInterceptor
	loggerStd.Printf("about init interceptor\n")
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = func(ctx context.Context) error {
			return nil
		}(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	loggerStd.Printf("success init interceptor\n")

	opts = append(opts, grpc.ConnectionTimeout(time.Duration(GRPCTIMEOUT)*time.Second))
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	//----------------------------------------------------
	if grpcOptInst.keepaliveTime != DEFKEEPALIVE {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     math.MaxInt64,
			MaxConnectionAge:      math.MaxInt64,
			MaxConnectionAgeGrace: math.MaxInt64,
			Time:                  grpcOptInst.keepaliveTime,
			Timeout:               grpcOptInst.keepaliveTimeout,
		}))
	}

	opts = append(opts, grpc.MaxConcurrentStreams(grpcOptInst.maxConcurrentStreams))
	opts = append(opts, grpc.MaxRecvMsgSize(grpcOptInst.maxReceiveMessageSize))
	opts = append(opts, grpc.MaxSendMsgSize(grpcOptInst.maxSendMessageSize))
	opts = append(opts, grpc.InitialWindowSize(grpcOptInst.initialWindowSize))
	opts = append(opts, grpc.InitialConnWindowSize(grpcOptInst.initialConnWindowSize))
	opts = append(opts, grpc.WriteBufferSize(grpcOptInst.writeBufferSize))
	opts = append(opts, grpc.ReadBufferSize(grpcOptInst.readBufferSize))
	//----------------------------------------------------
	loggerStd.Printf("about to call grpc.NewServer(opts...),maxRecv:%v,maxSend:%v\n", grpcOptInst.maxReceiveMessageSize, grpcOptInst.maxSendMessageSize)
	x.grpcserver = grpc.NewServer(opts...)
	addKillerCheck(killerNormalPriority, "grpcserver", x)

	loggerStd.Printf("about to call utils.RegisterXsfCallServer(x.grpcserver, srv)\n")

	utils.RegisterXsfCallServer(x.grpcserver, srv)
	utils.RegisterToolBoxServer(x.grpcserver, &ToolBoxServer{})
	loggerStd.Printf("about to call reflection.Register(x.grpcserver)\n")

	reflection.Register(x.grpcserver)

	loggerStd.Println("about to exec userCallback")
	dealUserCallBack()

	//----------------------------------------------------
	loggerStd.Println("about to call x.grpcserver.Serve")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := x.grpcserver.Serve(listener); err != nil {
			resErr = err
		}
	}()

	{
		//health check
		ctxTm, cancelTm := context.WithTimeout(context.Background(), grpcSleepTime)
		defer cancelTm()

		loggerStd.Printf("about to check if the grpc service(%v) is started\n", listener.Addr().String())
	retryEnd:
		for {
			select {
			case <-ctxTm.Done():
				{
					log.Panicf("deadlineExceed when check grpc server(%v)\n", listener.Addr().String())
				}
			default:
				{
					if healthCheck(listener.Addr().String()) {
						loggerStd.Printf("grpc server(%v) started successfully\n", listener.Addr().String())
						break retryEnd
					} else {
						loggerStd.Printf("about to reconfirm if the grpc service(%v) is started successfully\n", listener.Addr().String())
					}
				}
			}
		}
	}
	{
		//bvt check
		utils.SyncBvtInitStatus()
		if err := bvt.Check(); err != nil {
			loggerStd.Println("bvtCheck:", err)
			os.Exit(bvt.ExceptionBvt)
		}
	}
	//----------------------------------------------------

	loggerStd.Printf("about to call finderadapter.Register(%v)\n", listener.Addr().String())
	finderExist, finderRegisterErr := finderadapter.Register(listener.Addr().String(), bc.CfgData.ApiVersion)
	if finderRegisterErr != nil {
		log.Panicf("finderadapter.Register fail -> addr:%v,bc:%+v,finderRegistErr:%v\n", listener.Addr().String(), bc, finderRegisterErr)
	}

	if finderExist {
		loggerStd.Printf("finderadapter.Register success. -> addr:%v\n", listener.Addr().String())
	}

	loggerStd.Println("about to exec fcDelayInst")
	fcDelayInst.exec()

	{
		//保留启动时间
		globalStart = time.Now()
	}
	waitCtx, waitCtxCancel := context.WithCancel(context.Background())

	addKillerCheck(killerLastPriority, "WaitForExit", &killerWrapper{callback: func() {
		waitCtxCancel()
	}})

	loggerStd.Println("blocking for grpcserver.Serve")

	utils.SyncFinishStatus()

	wg.Wait()
	loggerStd.Println("success stop grpc graceful.")

	loggerStd.Println("Waiting for exit.")

	<-waitCtx.Done()

	loggerStd.Println("all mission complete.")

	return nil
}
