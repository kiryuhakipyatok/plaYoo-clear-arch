package config

import (
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.InfoLevel)
	logrus.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{})

	return log
}