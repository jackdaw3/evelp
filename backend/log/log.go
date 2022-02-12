package log

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func Init() {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
}

func SetLevel(level string) {
	switch level {
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	}
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(err error, a ...interface{}) {
	if logger.IsLevelEnabled(logrus.ErrorLevel) {
		fmt.Printf("\ntime=\"%v\" level=error msg=\"%v\"\n%+v\n\n", time.Now().Format(time.RFC3339), fmt.Sprint(a...), err)
	}
}

func Errorf(err error, format string, a ...interface{}) {
	if logger.IsLevelEnabled(logrus.ErrorLevel) {
		fmt.Printf("\ntime=\"%v\" level=error msg=\"%v\"\n%+v\n\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...), err)
	}
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
