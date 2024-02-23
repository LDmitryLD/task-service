package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type AppConf struct {
	AppName   string
	Server    Server
	RPCServer RPCServer
	UserRPC   UserRPC
	TaskRPC   TaskRPC
	Logger    Logger
}

type Logger struct {
	Level string
}

type Server struct {
	Port            string
	ShutdownTimeout time.Duration
}

type RPCServer struct {
	Port string
	Type string
}

type UserRPC struct {
	Host string
	Port string
}

type TaskRPC struct {
	Host string
	Port string
}

func NewAppConf() AppConf {
	port := os.Getenv("SERVER_PORT")

	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		Server: Server{
			Port: port,
		},
		Logger: Logger{
			Level: os.Getenv("LOGGER_LEVEL"),
		},
	}
}

func (a *AppConf) Init(logger *zap.Logger) {
	shutDownTimeOut, err := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT"))
	if err != nil {
		logger.Fatal("config: parse shutdown timeout error", zap.Error(err))
	}
	shutDownTimeout := time.Duration(shutDownTimeOut) * time.Second

	a.Server.ShutdownTimeout = shutDownTimeout

	a.RPCServer.Port = os.Getenv("RPC_PORT")
	a.RPCServer.Type = os.Getenv("RPC_PROTOCOL")

	a.UserRPC.Port = os.Getenv("USER_RPC_PORT")
	a.UserRPC.Host = os.Getenv("USER_RPC_HOST")

	a.TaskRPC.Port = os.Getenv("TASK_RPC_PORT")
	a.TaskRPC.Host = os.Getenv("TASK_RPC_HOST")

}
