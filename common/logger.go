package common

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Entry

func init() {
	std := logrus.StandardLogger()
	std.Level = logrus.DebugLevel
	std.SetFormatter(&logrus.JSONFormatter{})
	Logger = logrus.NewEntry(std)
}
