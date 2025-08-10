package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var log *logrus.Logger

func Init(level string) {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	lvl, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to 'info'", level)
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(lvl)
	}
}

func GetLogger() *logrus.Logger {
	if log == nil {
		Init("info") // Initialize with default info level if not already initialized
	}
	return log
}
