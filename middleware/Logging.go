package middleware

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

var Logger *logrus.Entry

func Log(error interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Logger.WithField("file", file).
		WithField("line", line).
		Info(error)
}
