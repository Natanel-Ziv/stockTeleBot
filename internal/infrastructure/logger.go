package infrastructure

import (
    "os"
    "github.com/sirupsen/logrus"
)


func NewLogger(logLevel string) *logrus.Logger {
    logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	logLvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Panic(err)
	}

	logger.SetLevel(logLvl)
	return logger
}