package log

import (
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logPath  string
	logLevel zapcore.Level //debug级别
	timeFmt  string        //格式化时间
	Core     *hertzzap.Logger
}

func GetLogger() *Logger {
	l := &Logger{
		logPath:  "./data/logs",
		logLevel: zap.InfoLevel,
		timeFmt:  "2006-01-02 15:04:05.000",
	}

	l.initZap()

	return l
}
