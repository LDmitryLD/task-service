package config

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type AppConf struct {
	AppName   string
	DB        DB
	RPCServer RPCServer
	TaksRPC   TaskRPC
	Logger    Logger
}

type Logger struct {
	Level string
}

type DB struct {
	Driver   string
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	MaxConn  int
	Timeout  int
}

type TaskRPC struct {
	Host string
	Port string
}

type RPCServer struct {
	Port string
	Type string
}

func NewAppConf() AppConf {

	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		DB: DB{
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Logger: Logger{
			Level: os.Getenv("LOGGER_LEVEL"),
		},
	}
}

func (a *AppConf) Init(logger *zap.Logger) {
	dbTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		logger.Fatal("config: parse db timeout error", zap.Error(err))
	}
	dbMaxConn, err := strconv.Atoi(os.Getenv("MAX_CONN"))
	if err != nil {
		logger.Fatal("config: parse db max connection error", zap.Error(err))
	}
	a.DB.Timeout = dbTimeout
	a.DB.MaxConn = dbMaxConn

	a.RPCServer.Port = os.Getenv("RPC_PORT")
	a.RPCServer.Type = os.Getenv("RPC_PROTOCOL")

	a.TaksRPC.Port = os.Getenv("TASK_RPC_PORT")
	a.TaksRPC.Host = os.Getenv("TASK_RPC_HOST")
}
