package conf

import (
	"flag"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/BurntSushi/toml"
	"fmt"
	"schemas"
	"common"
)

type (
	Config struct {
		Auth    Auth    `toml:"auth"`
		Xsf     Xsf     `toml:"xsf"`
		Session Session `toml:"session"`
		Log     Log     `toml:"log"`
		Schema  Schema  `toml:"schema"`
		Server  Server  `toml:"server"`
	}

	Auth struct {
		EnableAuth      bool   `toml:"enable_auth"`  // 废弃
		EnableAppidCheck bool  `toml:"enable_appid_check"` // 是否开启appid 校验
		SecretKey       string `toml:"secret_key"`
		MaxDateInterval int    `toml:"max_date_interval"`

	}
	Xsf struct {
		ServerPort   string `toml:"server_port"`   //xsf server端口
		CacheService bool   `toml:"cache_service"`
		CacheConfig  bool   `toml:"cache_config"`
		CachePath    string `toml:"cache_path"`
		CallRetry    int    `toml:"call_retry"`
		Location     string `toml:"location"`
		From         string `toml:"from"`
		EnableRespsort bool `toml:"enable_respsort"`
		Dc  string          `toml:"dc"`
		XsfLocalIp   string
		SpillEnable bool `toml:"spill_enable"`
	}

	Session struct {
		ScanInterver     int `toml:"scan_interver"`   // session 全局扫描间隔
		TimeoutInterver  int `toml:"timeout_interver"` // session 等待超时|连接等待超时
		HandshakeTimeout int `toml:"handshake_timeout"`  // 握手超时
		SessionCloseWait int `toml:"session_close_wait"` // session关闭等待时间
		SessionTimeout   int `toml:"session_timeout"` // session时长限制

	}


	Log struct {
		Level  string `toml:"level"`
		File   string `toml:"file"`
		Count  int    `toml:"count"`
		Size   int    `toml:"size"`
		Caller bool   `toml:"caller"`
		Batch  int    `toml:"batch"`
		Asyn   bool   `toml:"async"`
	}

	Schema struct {
		Enable  bool     `toml:"enable"`
		Services []string `toml:"services"`
	}

	Center struct {
		Project      string
		Group        string
		Service      string
		Version      string
		CompanionUrl string
	}

	Server struct {
		WriteFirst bool  `toml:"write_first"`
		Mock      bool   `toml:"mock"`
		Host      string `toml:"host"`
		Mode      string `toml:"mode"`
		NetCard   string `toml:"net_card"`
		Port      string `toml:"port"`
		ConsolLog string `toml:"consol_log"`
		PipeDepth int    `toml:"pipe_depth"`
		PipeTimeout int  `toml:"pipe_timeout"`
		EnableSonar bool `toml:"enable_sonar"`
		IgnoreRespCodes []int `toml:"ignore_resp_codes"`
		IgnoreSonarCodes []int `toml:"ignore_sonar_codes"`
		MaxConn int 	 `toml:"max_conn"`
		EnableConnLimit bool `toml:"enable_conn_limit"`
		AdminListen    string `toml:"admin_listen"`
		ScriptEnable bool `toml:"script_enable"`

	}
)

var (
	Conf   Config
	Centra Center
)

var (

	project  = flag.String("project", "guiderAllService", "project name")
	group    = flag.String("group", "gas", "group name")
	service  = flag.String("service", "webgate-ws", "service name ;should be same as  main tag in xsf.toml")
	version  = flag.String("version", "0.9.3_trace", "config center version")
	url      = flag.String("url", "http://10.1.87.70:6868", "config center companionUrl")
	cfg      = flag.String("cfg", "app.toml", "name of config file")
	schema      = flag.String("schema", "schemas/schema.json", "name of schema file ,this file is reqiured while bootMode is native ")
	BootMode = flag.Bool("nativeBoot", false, "boot from native config ")
)

const (

	APP_CONFIG = "app.toml"
	SCNEMA = "schema.json"
	LimitConf = "limit.json"

)

func InitConf() {
	flag.Parse()

	Centra = Center{
		Project:      *project,
		Group:        *group,
		Service:      *service,
		CompanionUrl: *url,
		Version:      *version,
	}
	//var f *finder.FinderManager
	if *BootMode {
		_, err := toml.DecodeFile(*cfg, &Conf)
		if err != nil {
			panic(err)
		}
		if err := schemas.LoadRouteMappingFromFile(*schema); err != nil {
			panic("cannot load mapping file:"+err.Error())
		}

		common.SetSystemEnv("APP_PORT",Conf.Server.Port)

	} else {
		InitCentra()
	}


	//获取本机ip
	ip, _ := utils.HostAdapter(Conf.Server.Host, Conf.Server.NetCard)
	Conf.Server.Host = ip
	Conf.Xsf.XsfLocalIp = ip + ":" + Conf.Xsf.ServerPort
	fmt.Printf("init conf:%+v\n", Conf)

	common.Setenv("APP_PORT",Conf.Server.Port)
	common.Setenv("APP_HOST",Conf.Server.Host)
	//InitSchema()
}

//获取端口号
func (server *Server) GetPort() string {
	return server.Port
}


func IsIgnoreRespCode(code int) bool {
	for i:=0;i< len(Conf.Server.IgnoreRespCodes);i++{
		if Conf.Server.IgnoreRespCodes[i] == code{
			return true
		}
	}
	return false
}

func IsIgnoreSonarCode(code int) bool {
	for i:=0;i< len(Conf.Server.IgnoreSonarCodes);i++{
		if Conf.Server.IgnoreSonarCodes[i] == code{
			return true
		}
	}
	return false
}
// 配置文件初始化
func (c *Config)Init()  {
	if c.Session.SessionTimeout<=0{
		c.Session.SessionTimeout = 30
	}

	if c.Session.HandshakeTimeout<=0{
		c.Session.HandshakeTimeout  = 4
	}

	if c.Session.TimeoutInterver <= 0{
		c.Session.TimeoutInterver = 15
	}

	if c.Session.SessionCloseWait<=0{
		c.Session.SessionCloseWait = 5
	}

	if c.Session.ScanInterver<=0{
		c.Session.ScanInterver = 30
	}
}