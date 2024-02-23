package server

import (
	"context"
	"fmt"
	"net"
	"projects/LDmitryLD/task-service/task/config"
	"projects/LDmitryLD/task-service/task/rpc/task"

	pb "github.com/LDmitryLD/protos2/gen/task"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server interface {
	Serve(ctx context.Context) error
}

type ServerGRPC struct {
	conf   config.RPCServer
	srv    *grpc.Server
	task   *task.TaskServiceGRPC
	logger *logrus.Logger
}

func NewServerGRPC(conf config.RPCServer, task *task.TaskServiceGRPC, logger *logrus.Logger) Server {
	return &ServerGRPC{
		conf:   conf,
		task:   task,
		srv:    grpc.NewServer(),
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
			s.logger.Error("gRPC server register error:", err.Error())
			chErr <- err
		}

		s.logger.Info("gRPC server started on:", l.Addr().String())

		pb.RegisterTaskerServer(s.srv, s.task)

		if err = s.srv.Serve(l); err != nil {
			s.logger.Error("gRPC server error:", err)
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
