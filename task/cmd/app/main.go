package main

import (
	"os"
	"projects/LDmitryLD/task-service/task/config"
	"projects/LDmitryLD/task-service/task/internal/infrastructure/logs"
	"projects/LDmitryLD/task-service/task/run"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	conf := config.NewAppConf()

	logger := logs.NewLogger(conf)

	conf.Init(logger)

	app := run.NewApp(conf, logger)

	if err := app.Bootstrap().Run(); err != nil {
		logger.WithField("err", err).Error("app run error")
		os.Exit(2)
	}
}
