package helpers

import (
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var logger *logrus.Logger

func LoggerInit() {
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "02-01-2006 15:04:05",
			LogFormat:       "[%lvl%] [%component%] %time% - %msg% \n",
		},
	}

	l := &lumberjack.Logger{
		Filename:   LogsPath + "/latest.log",
		MaxSize:    1, // MegaBytes
		MaxBackups: 8, // Max Files
		MaxAge:     7, // Days
		Compress:   false,
	}
	mWriter := io.MultiWriter(os.Stdout, l)
	logger.SetOutput(mWriter)
}

func LogError(message string, component string) {
	logger.WithField("component", component).Error(message)
}

func LogInfo(message string, component string) {
	logger.WithField("component", component).Info(message)
}

func LogWarning(message string, component string) {
	logger.WithField("component", component).Warningln(message)
}
