package server

import (
	"context"
	"net/http"
	"projects/LDmitryLD/task-service/gateway/config"

	"go.uber.org/zap"
)

type Server interface {
	Serve(ctx context.Context) error
}

type HTTPServer struct {
	conf config.Server
	srv  *http.Server
	log  *zap.Logger
}

func NewHTTPServer(conf config.Server, server *http.Server, logger *zap.Logger) Server {
	return &HTTPServer{conf: conf, srv: server, log: logger}
}

func (s *HTTPServer) Serve(ctx context.Context) error {
	var err error

	chErr := make(chan error)

	go func() {
		s.log.Info("server started", zap.String("addr", s.srv.Addr))
		if err = s.srv.ListenAndServe(); err != http.ErrServerClosed {
			s.log.Error("http listen and serve error", zap.Error(err))
			chErr <- err
		}
	}()

	select {
	case <-chErr:
		return err
	case <-ctx.Done():
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), s.conf.ShutdownTimeout)
	defer cancel()
	err = s.srv.Shutdown(ctxShutdown)

	return err
}
