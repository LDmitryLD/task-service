package run

import (
	"context"
	"os"
	"projects/LDmitryLD/task-service/task/config"
	"projects/LDmitryLD/task-service/task/internal/db"
	"projects/LDmitryLD/task-service/task/internal/infrastructure/server"
	"projects/LDmitryLD/task-service/task/internal/modules/task/service"
	"projects/LDmitryLD/task-service/task/internal/modules/task/storage"
	userservice "projects/LDmitryLD/task-service/task/internal/modules/user/service"
	"projects/LDmitryLD/task-service/task/rpc/task"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run() error
}

type App struct {
	conf   config.AppConf
	logger *logrus.Logger
	rpc    server.Server
	Sig    chan os.Signal
}

func NewApp(conf config.AppConf, logger *logrus.Logger) *App {
	return &App{
		conf:   conf,
		logger: logger,
		Sig:    make(chan os.Signal, 1),
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.WithField("os_signal", sigInt).Info("signal interrupt recieved")
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.rpc.Serve(ctx)
		if err != nil {
			a.logger.Error("app: server error:", err)
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) Bootstrap(options ...interface{}) Runner {
	_, sqlAdapter, err := db.NewSqlDB(a.conf.DB, a.logger)
	if err != nil {
		a.logger.WithField("err", err).Fatal("error init db")
	}

	userClient, err := userservice.NewUserGRPCClient(a.conf.UserRPC)
	if err != nil {
		a.logger.WithField("err", err).Fatal("error init user rpc client")
	}

	taskStorage := storage.NewTaskStorage(sqlAdapter)

	taskService := service.NewTasker(taskStorage, userClient, a.logger)

	taskServiceGRPC := task.NewTaskServiceGRPC(taskService)

	a.rpc = server.NewServerGRPC(a.conf.RPCServer, taskServiceGRPC, a.logger)

	return a
}
