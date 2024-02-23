package main

import (
	"os"
	"projects/LDmitryLD/task-service/user/config"
	"projects/LDmitryLD/task-service/user/internal/infrastructure/logs"
	"projects/LDmitryLD/task-service/user/run"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()

	conf := config.NewAppConf()

	logger := logs.NewLogger(conf, os.Stdout)

	conf.Init(logger)

	app := run.NewApp(conf, logger)

	if err := app.Bootstrap().Run(); err != nil {
		logger.Error("app run error", zap.Error(err))
		os.Exit(2)
	}
}
