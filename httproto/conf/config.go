package conf

import (
	"sync/atomic"
	"unsafe"
)

type (
	Config struct {
		Session Session `toml:"session"`
		Server  Server  `toml:"server"`
		Log     Log     `toml:"log"`
	}
	Session struct {
		ScanInterver     int `toml:"scan_interver"`      // session 全局扫描间隔
		TimeoutInterver  int `toml:"timeout_interver"`   // session 等待超时|连接等待超时
		HandshakeTimeout int `toml:"handshake_timeout"`  // 握手超时
		SessionCloseWait int `toml:"session_close_wait"` // session关闭等待时间
		SessionTimeout   int `toml:"session_timeout"`    // session时长限制
		ReadTimeout      int `toml:"read_timeout"`
	}
	Server struct {
		WriteFirst      bool   `toml:"write_first"`       //write response wherever first frame has result
		Mock            bool   `toml:"mock"`              // deprecated ,use schema.mapping.mock instead
		Host            string `toml:"host"`              // listen host ,if empty, server will listen at first net card
		Mode            string `toml:"mode"`              // release or debug
		NetCard         string `toml:"net_card"`          //
		Port            string `toml:"port"`              // listen port
		PipeDepth       int    `toml:"pipe_depth"`        // deprecated ,
		PipeTimeout     int    `toml:"pipe_timeout"`      // deprecated
		EnableSonar     bool   `toml:"enable_sonar"`      // true if use sonar log
		IgnoreRespCodes []int  `toml:"ignore_resp_codes"` // code return from upstream will no longer return to client
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
)

func (c *Config) Init() {
	if c.Session.SessionTimeout <= 0 {
		c.Session.SessionTimeout = 150
	}

	if c.Session.HandshakeTimeout <= 0 {
		c.Session.HandshakeTimeout = 4
	}

	if c.Session.TimeoutInterver <= 0 {
		c.Session.TimeoutInterver = 15
	}

	if c.Session.SessionCloseWait <= 0 {
		c.Session.SessionCloseWait = 5
	}

	if c.Session.ScanInterver <= 0 {
		c.Session.ScanInterver = 30
	}
	if c.Session.ReadTimeout <= 0 {
		c.Session.ReadTimeout = 10
	}
	c.Server.IgnoreRespCodes = []int{10101}
	c.Log.Level = "info"

}

// config instance 并发安全
var confInstance unsafe.Pointer

func GetConfInstance() *Config {
	return (*Config)(atomic.LoadPointer(&confInstance))
}

func InitConfig() error {
	c := &Config{}
	c.Init()

	atomic.StorePointer(&confInstance, unsafe.Pointer(c))
	return nil
}
