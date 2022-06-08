package log

import (
	"log"
	"os"
)

var Log Logger

type Logger interface {
	Infof(fmt string, v ...interface{})
	Debugf(fmt string, v ...interface{})
	Errorf(fmt string, v ...interface{})
	Printf(fmt string, v ...interface{})
}

type DefaultLogger struct {
	defaultLog *defaultLogger
}

func NewDefaultLogger() Logger {

	logFile, err := os.OpenFile("findergo.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("create log err ：", err)
	}
	SetPrefix("findergo  ---- ")
	SetOutput(logFile)
	SetFlags(Lshortfile | Lmicroseconds | Ldate)

	logger := &DefaultLogger{defaultLog: defaultStd}
	return logger
}

func (l *DefaultLogger) Info(v ...interface{}) {
	l.defaultLog.Println(v)
}

func (l *DefaultLogger) Debug(v ...interface{}) {
	l.defaultLog.Println(v)
}
func (l *DefaultLogger) Printf(fmt string, v ...interface{}) {
	l.defaultLog.Printf(fmt, v)
}
func (l *DefaultLogger) Error(v ...interface{}) {
	l.defaultLog.Println(v)
}

func (l *DefaultLogger) Infof(fmt string, v ...interface{}) {
	l.defaultLog.Printf(fmt, v)
}

func (l *DefaultLogger) Debugf(fmt string, v ...interface{}) {
	l.defaultLog.Printf(fmt, v)
}

func (l *DefaultLogger) Errorf(fmt string, v ...interface{}) {
	l.defaultLog.Printf(fmt, v)
}
