package xsf

import "git.xfyun.cn/AIaaS/xsf-external/utils"

type options struct {
	router       *OpRouter
	rateFallback RateInterface
}
type ServerOption func(*options)

func SetOpRouter(router *OpRouter) ServerOption {
	return func(o *options) {
		o.router = router
	}
}

func SetRateFallback(rateFallback RateInterface) ServerOption {
	return func(o *options) {
		o.rateFallback = rateFallback
	}
}

type RateInterface interface {
	Call(*Req, *Span) (*Res, error)
}
type UserInterface interface {
	Init(*ToolBox) error
	FInit() error
	Call(*Req, *Span) (*Res, error)
}

type XsfServer struct {
	Toolbox      *ToolBox
	SidGenerator *XrpcSidGenerator
	UserImpl     UserInterface
	Opts         *options
	Bc           BootConfig
}

func (x *XsfServer) init(bc BootConfig, srv UserInterface, opt ...ServerOption) error {
	x.UserImpl = srv
	x.Bc = bc
	var toolbox ToolBox
	var err error
	if err = toolbox.Init(bc); nil != err {
		return err
	}
	x.Opts = &options{}
	for _, o := range opt {
		o(x.Opts)
	}
	x.Toolbox = &toolbox
	x.Toolbox.Log.Infow("toolbox init complete.")
	return err
}
func (x *XsfServer) run() error {
	return xrpcsRunWrapper(x.Bc, x.Toolbox, x.UserImpl, x.Opts)
}
func (x *XsfServer) Run(bc BootConfig, srv UserInterface, opt ...ServerOption) error {
	bc, bcErr, bcExt := bcCheck(bc)
	//bcConfig配置不合法
	if nil != bcErr {
		return bcErr
	}
	if utils.Native == bc.CfgMode && "" != bc.CfgData.CfgName {
		bc.CfgData.CfgName = utils.FileNamePreProcessing(bc.CfgData.CfgName)
	}
	if err := x.init(bc, srv, opt...); nil != err {
		return err
	}
	//打印bcConfig检查中的警告
	x.Toolbox.Log.Warnf("bcCheck:%v\n", bcExt)
	return x.run()
}
