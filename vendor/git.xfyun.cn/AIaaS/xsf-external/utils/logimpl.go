package utils

import (
	//"genitus/flange"
	"log"
	"os"
)

var loggerStd = (&LoggerStderr{}).Init("")

type LoggerStderr struct {
	*log.Logger
}

func (l *LoggerStderr) Init(prefix string) *LoggerStderr {
	l.Logger = log.New(os.Stderr, prefix, log.LstdFlags)
	return l
}

type customLogImpl struct {
	//flange.CustomLogInterface
	log *Logger
}

func NewLogImpl(l *Logger) *customLogImpl {
	cli := new(customLogImpl)
	cli.log = l
	return cli
}

func (cli *customLogImpl) Infof(format string, params ...interface{}) {
	cli.log.Infof(format, params...)
}

func (cli *customLogImpl) Debugf(format string, params ...interface{}) {
	cli.log.Debugf(format, params...)
}

func (cli *customLogImpl) Errorf(format string, params ...interface{}) {
	cli.log.Errorf(format, params...)
}

func (cli *customLogImpl) Info(params ...interface{}) {
	cli.log.Infof("", params...)
}

func (cli *customLogImpl) Debug(params ...interface{}) {
	cli.log.Debugf("", params...)
}

func (cli *customLogImpl) Error(params ...interface{}) {
	cli.log.Errorf("", params...)
}

var dbgLoggerStd = newDbsLoggerStd("xxx ", false)

func init() {
	xsfDbg := os.Getenv("XSF-DEBUG")
	if "1" == xsfDbg {
		dbgLoggerStd = newDbsLoggerStd("debug=>", true)
	}
}

type DbgLoggerStd struct {
	logger *LoggerStderr
	able   bool
}

func newDbsLoggerStd(prefix string, able bool) *DbgLoggerStd {
	l := DbgLoggerStd{
		logger: (&LoggerStderr{}).Init(prefix),
		able:   able,
	}
	return &l
}
func (d *DbgLoggerStd) Printf(format string, v ...interface{}) {
	if !d.able {
		return
	}
	d.logger.Printf(format, v...)
}
func (d *DbgLoggerStd) Println(v ...interface{}) {
	if !d.able {
		return
	}
	d.logger.Println(v...)
}
