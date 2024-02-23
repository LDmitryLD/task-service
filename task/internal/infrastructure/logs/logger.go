package logs

import (
	"projects/LDmitryLD/task-service/task/config"

	"github.com/sirupsen/logrus"
)

func NewLogger(conf config.AppConf) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(conf.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		ForceColors:   true,
	})

	return logger
}
