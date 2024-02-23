package run

import (
	"context"
	"os"
	"projects/LDmitryLD/task-service/user/config"
	"projects/LDmitryLD/task-service/user/internal/db"
	"projects/LDmitryLD/task-service/user/internal/infrastructure/server"
	taskservice "projects/LDmitryLD/task-service/user/internal/modules/task/service"
	"projects/LDmitryLD/task-service/user/internal/modules/user/service"
	"projects/LDmitryLD/task-service/user/internal/modules/user/storage"
	"projects/LDmitryLD/task-service/user/rpc/user"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run() error
}

type App struct {
	conf   config.AppConf
	logger *zap.Logger
	rpc    server.Server
	Sig    chan os.Signal
}

func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sigInt))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.rpc.Serve(ctx)
		if err != nil {
			a.logger.Error("app: server error", zap.Error(err))
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
		a.logger.Fatal("error init db", zap.Error(err))
	}

	//components := component.NewComponents(a.conf, a.logger)

	taskClient, err := taskservice.NewTaskGRPCClient(a.conf.TaksRPC)
	if err != nil {
		a.logger.Fatal("error init task rpc client", zap.Error(err))
	}

	userStorage := storage.NewUserStorage(sqlAdapter)

	userService := service.NewUser(userStorage, taskClient, a.logger)

	userServiceGRPC := user.NewUserServiceGRPC(userService)

	a.rpc = server.NewServerGRPC(a.conf.RPCServer, userServiceGRPC, a.logger)

	return a
}
