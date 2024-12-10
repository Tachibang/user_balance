package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetLogrus(level string) *logrus.Logger {
	logger := logrus.New()

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Fatalf("Неверный уровень логирования: %v", err)
	}
	logger.SetLevel(logLevel)

	logger.SetOutput(os.Stdout)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}
