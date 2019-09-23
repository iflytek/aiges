package client

import (
	"config"
	"consts"
	"util"

	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

var rpcClient *RpcClient

//rpc client结构体
type RpcClient struct {
	client  *xsf.Client
	timeout int
}

//初始化rpc client
func InitRpcClient() {
	cfgName := consts.XSFC_FILE
	if config.UseCfgCentre == 0 {
		cfgName = consts.PREFIX + consts.XSFC_FILE
	}
	//定义相关的启动参数
	cli, err := xsf.InitClient(config.Service,
		utils.CfgMode(config.UseCfgCentre),
		utils.WithCfgName(cfgName),
		utils.WithCfgURL(config.CompanionUrl),
		utils.WithCfgPrj(config.Project),
		utils.WithCfgGroup(config.Group),
		utils.WithCfgService(config.Service),
		utils.WithCfgVersion(config.Version),
	)
	if err != nil {
		util.SugarLog.Errorw("InitRpcClient failed", "err", err)
		panic(err)
	}

	//读取业务所需要的若干配置
	cfg := cli.Cfg()
	tm, e := cfg.GetInt("sfc", "timeout")
	if e != nil {
		tm = 3000
	}

	//创建RpcClient
	rpcClient = &RpcClient{client: cli, timeout: tm}
}

//获取rpc client
func GetRpcClient() *RpcClient {
	return rpcClient
}

//创建caller可全局使用
func (c *RpcClient) CreateCaller() *xsf.Caller {
	return xsf.NewCaller(c.client)
}
