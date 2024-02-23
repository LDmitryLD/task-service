package main

import (
	"os"
	"projects/LDmitryLD/task-service/gateway/config"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/logs"
	"projects/LDmitryLD/task-service/gateway/run"

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
