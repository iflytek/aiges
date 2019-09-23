package utils

import (
	"testing"
	"time"
)

func TestNewCentreWithFinder(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	co := &CfgOption{
		name:         "lbv2.toml",
		tick:         time.Second,
		stmout:       time.Second,
		url:          "http://10.1.87.69:6868",
		cachePath:    ".",
		cacheConfig:  true,
		cacheService: true,
		prj:          "guiderAllService",
		group:        "gas",
		srv:          "lbv2",
		ver:          "2.2.7",
		log: func() *Logger {
			logger, err := NewLocalLog(
				SetCaller(true),
				SetLevel("debug"),
				SetFileName("test.log"),
				SetMaxSize(3),
				SetMaxBackups(3),
				SetMaxAge(3),
				SetAsync(false),
				SetCacheMaxCount(30000),
				SetBatchSize(1024))
			checkErr(err)
			return logger
		}()}
	findInst, findInstErr := NewFinder(co)
	checkErr(findInstErr)
	co.fm = findInst

	cfg, cfgErr := NewCentreWithFinder(co)
	checkErr(cfgErr)

	loggerStd.Println(cfg.GetInt("bo", "ticker"))
	time.Sleep(time.Second * 10)
	loggerStd.Println(cfg.GetInt("bo", "ticker"))

}
