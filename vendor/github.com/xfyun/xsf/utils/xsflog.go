/*
* @file	xsflog.go
* @brief	simple zap adapter (refers to official website)
* @author	sqjian
* @version	1.0
* @date		2017.11.29
 */
package utils

import (
	"errors"
	lumberjack "github.com/xfyun/lumberjack-ccr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var globalpid int64

func init() {
	globalpid = int64(os.Getpid())
}

type xsflogger = zap.SugaredLogger

var async = true
var cacheMaxCount = -1
var batchSize = 16 * 1024
var defaultWash = 60

type Logger struct {
	lumberjack *lumberjack.Logger
	logger     *xsflogger
}
type loggerCfg struct {
	logLevel      string
	fileName      string
	maxSize       int
	maxBackups    int
	maxAge        int
	async         bool
	cacheMaxCount int
	batchSize     int
	wash          int
	caller        bool
}
type optionFunc func(*loggerCfg)
type Option interface {
	apply(*loggerCfg)
}

func (f optionFunc) apply(l *loggerCfg) {
	f(l)
}
func SetCaller(Caller bool) Option {
	return optionFunc(func(l *loggerCfg) {
		l.caller = Caller
	})
}
func SetAsync(async bool) Option {
	return optionFunc(func(l *loggerCfg) {
		l.async = async
	})
}
func SetLevel(logLevel string) Option {
	return optionFunc(func(l *loggerCfg) {
		l.logLevel = logLevel
	})
}
func SetFileName(fileName string) Option {
	return optionFunc(func(l *loggerCfg) {
		l.fileName = fileName
	})
}
func SetMaxSize(maxSize int) Option {
	return optionFunc(func(l *loggerCfg) {
		l.maxSize = maxSize
	})
}
func SetMaxBackups(maxBackups int) Option {
	return optionFunc(func(l *loggerCfg) {
		l.maxBackups = maxBackups
	})
}
func SetMaxAge(maxAge int) Option {
	return optionFunc(func(l *loggerCfg) {
		l.maxAge = maxAge
	})
}

//缓存大小，单位条数,超过会丢弃，为-1时堆积数据至内存中
func SetCacheMaxCount(cacheMaxCount int) Option {
	return optionFunc(func(l *loggerCfg) {
		l.cacheMaxCount = cacheMaxCount
	})
}

//批处理大小，单位条数，一次写入条数（触发写事件的条数）
func SetBatchSize(batchSize int) Option {
	return optionFunc(func(l *loggerCfg) {
		l.batchSize = batchSize
	})
}

func SetWash(w int) Option {
	return optionFunc(func(l *loggerCfg) {
		l.wash = w
	})
}
func (l *Logger) Closeout() {
	loggerStd.Println("about to close lumberjack.")
	l.lumberjack.Stop()
}
func (l *Logger) Printf(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}
func (l *Logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}
func (l *Logger) Start(msg string, keysAndValues ...interface{}) {
	l.Errorw(" iFlyTEK log file start")
}

//func newLocalLog(logLevel string, fileName string, maxSize int, maxBackups int, maxAge int, async bool, cacheMaxCount int, batchSize int, caller bool) (*Logger, error) {
func newLocalLog(lc *loggerCfg) (*Logger, error) {
	lc.logLevel = strings.ToLower(lc.logLevel)
	if paramsCK := func() error {
		if lc.logLevel != "info" && lc.logLevel != "debug" && lc.logLevel != "warn" && lc.logLevel != "error" && lc.logLevel != "none" {
			return errors.New("params is illegal")
		} else {
			return nil
		}
	}(); paramsCK != nil {
		return nil, paramsCK
	}
	userPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if lc.logLevel == "none" {
			return false
		}
		return lvl >= func() zapcore.Level {
			switch lc.logLevel {
			case "debug":
				{
					return zapcore.DebugLevel
				}
			case "info":
				{
					return zapcore.InfoLevel
				}
			case "warn":
				{
					return zapcore.WarnLevel
				}
			case "error":
				{
					return zapcore.ErrorLevel
				}
			default:
				{
					return zapcore.ErrorLevel
				}
			}
		}()
	})
	lumberjackInst := &lumberjack.Logger{
		Filename:   lc.fileName,
		MaxSize:    lc.maxSize, // megabytes
		MaxBackups: lc.maxBackups,
		MaxAge:     lc.maxAge, // days

		Async:         lc.async,
		CacheMaxCount: lc.cacheMaxCount,
		BatchSize:     lc.batchSize,
		Wash:          lc.wash,
	}
	lumberjackInst.Start()
	logRotateUserWriter := zapcore.AddSync(lumberjackInst)
	commonEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	commonEncoder := zapcore.NewJSONEncoder(commonEncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(commonEncoder, logRotateUserWriter, userPriority),
	)
	l := &Logger{}
	if lc.caller {
		l = &Logger{lumberjack: lumberjackInst, logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Fields(zapcore.Field{Key: "pid", Type: zapcore.Int64Type, Integer: globalpid})).Sugar()}

	} else {
		l = &Logger{lumberjack: lumberjackInst, logger: zap.New(core, zap.AddCallerSkip(1), zap.Fields(zapcore.Field{Key: "pid", Type: zapcore.Int64Type, Integer: globalpid})).Sugar()}
	}
	return l, nil
}
func NewLocalLog(opt ...Option) (*Logger, error) {
	lc := &loggerCfg{async: async, cacheMaxCount: cacheMaxCount, batchSize: batchSize, wash: defaultWash}
	for _, o := range opt {
		o.apply(lc)
	}
	return newLocalLog(lc)
	//	return newLocalLog(loggercfg.logLevel, loggercfg.fileName, loggercfg.maxSize, loggercfg.maxBackups, loggercfg.maxAge, loggercfg.async, loggercfg.cacheMaxCount, loggercfg.batchSize, loggercfg.caller)
}
func StopLocalLog(logger *Logger) {
	logger.lumberjack.Stop()
}
