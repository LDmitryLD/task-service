package run

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"projects/LDmitryLD/task-service/gateway/config"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/component"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/logs/responder"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/router"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/server"
	"projects/LDmitryLD/task-service/gateway/internal/modules"
	taskservice "projects/LDmitryLD/task-service/gateway/internal/modules/task/service"
	userservice "projects/LDmitryLD/task-service/gateway/internal/modules/user/service"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run() error
}

type App struct {
	conf   config.AppConf
	logger *zap.Logger
	srv    server.Server
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
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Info("app: server error:", zap.Error(err))
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
	responseManager := responder.NewResponder(a.logger)

	components := component.NewComponents(a.conf, responseManager, a.logger)

	userClientRPC, err := userservice.NewUserGRPCClient(a.conf.UserRPC)
	if err != nil {
		a.logger.Fatal("error init user rpc client", zap.Error(err))
	}

	taskClientRPC, err := taskservice.NewTaskGRPCClient(a.conf.TaskRPC)
	if err != nil {
		a.logger.Fatal("error init task rpc client", zap.Error(err))
	}

	controllers := modules.NewControllers(userClientRPC, taskClientRPC, components)

	r := router.NewRouter(controllers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.conf.Server.Port),
		Handler: r,
	}

	a.srv = server.NewHTTPServer(a.conf.Server, srv, a.logger)

	return a
}
