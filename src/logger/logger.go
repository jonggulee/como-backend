package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

func init() {
	Logger = logrus.New()

	// JSONFormatter
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.DebugLevel)
}

func Debugf(reqId string, format string, v ...interface{}) {
	loggerWithField := Logger.WithFields(logrus.Fields{})
	loggerWithField.Data["file"] = fileInfo(2)

	loggerWithField.Data["requestId"] = reqId

	loggerWithField.Debugf(format, v...)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
