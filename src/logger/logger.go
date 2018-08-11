package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func New() *logrus.Logger {
	log := logrus.New()
	formatter := new(prefixed.TextFormatter)
	formatter.ForceFormatting = true
	log.Formatter = formatter

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "debug" {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}

	return log
}
