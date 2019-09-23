package utils

import (
	"testing"
)

func TestXsflogCaller(t *testing.T) {
	//logger, err := NewLocalLog("info", "test.log", 3, 3, 3)
	logger, err := NewLocalLog(SetCaller(true), SetLevel("error"), SetFileName("test.log"), SetMaxSize(3), SetMaxBackups(3), SetMaxAge(3), SetAsync(true), SetCacheMaxCount(30000), SetBatchSize(1024))
	if err != nil {
		t.Fatal(err)
	}
	defer StopLocalLog(logger)
	logger.Infof("just a test.")
	logger.Infof("just a test K:%v.", "this is k")
	logger.Debugf("just a test K:%v.", "this is k")
	logger.Errorf("just a test K:%v.", "this is k")
	logger.Debugw("this is log.", "k", "v")
	logger.Infow("this is log.", "k", "v")
	logger.Warnw("this is log.", "k", "v")
	logger.Errorw("this is log.", "k", "v")
}
func TestXsflogr(t *testing.T) {
	//logger, err := NewLocalLog("info", "test.log", 3, 3, 3)
	logger, err := NewLocalLog(SetCaller(false), SetLevel("info"), SetFileName("test.log"), SetMaxSize(3), SetMaxBackups(3), SetMaxAge(3), SetAsync(true), SetCacheMaxCount(30000), SetBatchSize(1024))
	if err != nil {
		t.Fatal(err)
	}
	defer StopLocalLog(logger)
	logger.Infof("just a test.")
	logger.Infof("just a test K:%v.", "this is k")
	logger.Debugf("just a test K:%v.", "this is k")
	logger.Errorf("just a test K:%v.", "this is k")
	logger.Debugw("this is log.", "k", "v")
	logger.Infow("this is log.", "k", "v")
	logger.Warnw("this is log.", "k", "v")
	logger.Errorw("this is log.", "k", "v")
}
func TestPerf(t *testing.T) {
	logger, err := NewLocalLog(SetCaller(false), SetLevel("info"), SetFileName("test.log"), SetMaxSize(3), SetMaxBackups(3), SetMaxAge(3), SetAsync(true), SetCacheMaxCount(30000), SetBatchSize(1024))
	if err != nil {
		t.Fatal(err)
	}
	defer StopLocalLog(logger)
	logger.Infof("just a test.")
	logger.Infof("just a test K:%v.", "this is k")
	logger.Debugf("just a test K:%v.", "this is k")
	logger.Errorf("just a test K:%v.", "this is k")
	logger.Debugw("this is log.", "k", "v")
	logger.Infow("this is log.", "k", "v")
	logger.Warnw("this is log.", "k", "v")
	logger.Errorw("this is log.", "k", "v")
	logger.Errorw("this is log.", "k", []string{"111", "222"})
}
