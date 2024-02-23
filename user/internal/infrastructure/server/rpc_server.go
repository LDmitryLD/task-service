package server

import (
	"context"
	"fmt"
	"net"
	"projects/LDmitryLD/task-service/user/config"
	"projects/LDmitryLD/task-service/user/rpc/user"

	pb "github.com/LDmitryLD/protos2/gen/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server interface {
	Serve(ctx context.Context) error
}

type ServerGRPC struct {
	conf   config.RPCServer
	srv    *grpc.Server
	user   *user.UserServiceGRPC
	logger *zap.Logger
}

func NewServerGRPC(conf config.RPCServer, user *user.UserServiceGRPC, logger *zap.Logger) Server {
	return &ServerGRPC{
		conf:   conf,
		srv:    grpc.NewServer(),
		user:   user,
		logger: logger,
	}
}

func (s *ServerGRPC) Serve(ctx context.Context) error {
	var err error

	chErr := make(chan error)
	go func() {
		var l net.Listener
		l, err = net.Listen("tcp", fmt.Sprintf(":%s", s.conf.Port))
		if err != nil {
			s.logger.Error("gRPC server register error:", zap.Error(err))
			chErr <- err
		}

		s.logger.Info("gRPC server started", zap.String("addr", l.Addr().String()))

		pb.RegisterUsererServer(s.srv, s.user)

		if err = s.srv.Serve(l); err != nil {
			s.logger.Error("gRPC server error", zap.Error(err))
			chErr <- err
		}
	}()

	select {
	case <-chErr:
		return err
	case <-ctx.Done():
		s.srv.GracefulStop()
	}
	return err
}
