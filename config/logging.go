package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func ConfigureLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
