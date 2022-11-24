package logger

import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	return l
}
