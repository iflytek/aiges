package server

import (
	"strconv"

	"config"
	"consts"
	"core"
	"util"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"time"
	"sync"

	"log"
)

//用户业务逻辑接口
type server struct {
	tool *xsf.ToolBox
}

//业务初始化接口
func (c *server) Init(toolbox *xsf.ToolBox) error {
	c.tool = toolbox
	config.ServiceHost = toolbox.NetManager.GetIp()
	config.ServicePort = strconv.Itoa(toolbox.NetManager.GetPort())
	return nil
}

//业务逆初始化接口
func (c *server) FInit() error {
	for i := 1; i <= 16; i++ {
		time.Sleep(1000 * time.Millisecond)
		util.SugarLog.Errorw("Finit", "Sleep", i)
	}

	return nil
}

//业务服务接口
func (c *server) Call(in *xsf.Req, span *xsf.Span) (*utils.Res, error) {
	res := xsf.NewRes()

	//检测数据
	if len(in.Data()) == 0 {
		//数据为空,输出日志
		util.SugarLog.Errorw("Call failed", "err", consts.ERR_MSG_BAD_RPC, "op", in.Op())
		//设置up result
		content := core.SetUpResult(nil, nil, nil, consts.MSP_ERROR_NO_DATA, consts.MSG_BAD_RPC)

		//返回结果
		if content != nil {
			data := xsf.NewData()
			data.Append(content)
			res.AppendData(data)
		}
		return res, nil
	}

	content := core.Core.Process(in.Op(), in.Data()[0].Data(), span, in.GetAllParam())

	//返回结果
	if content == nil {
		util.SugarLog.Errorw("response failed", "content is nil", consts.ERR_MSG_BAD_RPC, "op", in.Op())
		content = core.SetUpResult(nil, nil, nil, consts.MSP_ERROR_NO_DATA, consts.MSG_BAD_RPC)
	}
	data := xsf.NewData()
	data.Append(content)
	res.AppendData(data)
	return res, nil
}

//初始化rpc server
func InitRpcServer() {
	//定义一个服务实例
	var serverInst xsf.XsfServer
	//定义相关的启动参数
	cfgName := consts.XSFS_FILE
	if config.UseCfgCentre == 0 {
		cfgName = consts.PREFIX + consts.XSFS_FILE
	}
	bc := xsf.BootConfig{
		CfgMode: utils.CfgMode(config.UseCfgCentre),
		CfgData: xsf.CfgMeta{
			CfgName:      cfgName,
			Project:      config.Project,
			Group:        config.Group,
			Service:      config.Service,
			Version:      config.Version,
			CompanionUrl: config.CompanionUrl,
		},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			if err := recover(); err != nil {
				util.SugarLog.Errorw("InitRpcServer failed", "err", err)
			}
		}()
		/*
			1、启动服务
			2、若有异常直接报错，注意需用户自己实现协程等待
		*/
		if err := serverInst.Run(bc, &server{}); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}
