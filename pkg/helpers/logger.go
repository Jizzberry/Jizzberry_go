package helpers

import (
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"runtime"
	"strings"
)

var logger *logrus.Logger

func LoggerInit() {
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "02-01-2006 15:04:05",
			LogFormat:       "[%lvl%] [%component%] %time% - %msg%\n",
		},
	}

	l := &lumberjack.Logger{
		Filename:   LogsPath + "/latest.log",
		MaxSize:    1, // MegaBytes
		MaxBackups: 8, // Max Files
		MaxAge:     7, // Days
		Compress:   true,
	}
	mWriter := io.MultiWriter(os.Stdout, l)
	logger.SetOutput(mWriter)
}

func LogError(message string) {
	logger.WithFields(logrus.Fields{
		"component": func() string {
			pc, _, _, ok := runtime.Caller(2)
			return getCalledInfo(pc, ok)
		}(),
	}).Error(message)
}

func LogInfo(message string) {
	logger.WithFields(logrus.Fields{
		"component": func() string {
			pc, _, _, ok := runtime.Caller(2)
			return getCalledInfo(pc, ok)
		}(),
	}).Info(message)
}

func LogWarning(message string) {
	logger.WithFields(logrus.Fields{
		"component": func() string {
			pc, _, _, ok := runtime.Caller(2)
			return getCalledInfo(pc, ok)
		}(),
	}).Warningln(message)
}

func getCalledInfo(pc uintptr, ok bool) string {
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		split := strings.Split(details.Name(), "/")
		return split[len(split)-1]
	}
	return "Unknown"
}
