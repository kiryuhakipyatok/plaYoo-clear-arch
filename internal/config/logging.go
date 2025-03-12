package config

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{})
	file,err:=os.OpenFile("../../files/logging/logrus.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err!=nil{
		log.WithError(err).Info("failed to open log file, using default stdout")
	}
	log.SetOutput(io.MultiWriter(file,os.Stdout))
	return log
}
