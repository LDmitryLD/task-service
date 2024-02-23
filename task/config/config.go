package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type AppConf struct {
	AppName   string
	DB        DB
	RPCServer RPCServer
	UserRPC   UserRPC
	Logger    Logger
}

type Logger struct {
	Level string
}

type DB struct {
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	MaxConn  int
	Timeout  int
}

type RPCServer struct {
	Port string
	Type string
}

type UserRPC struct {
	Host string
	Port string
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

func (a *AppConf) Init(logger *logrus.Logger) {
	dbTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Fatal("config: parse db timeout error")
	}
	dbMaxConn, err := strconv.Atoi(os.Getenv("MAX_CONN"))
	if err != nil {
		logger.WithField("error", err).Fatal("config: parse db max conn error")
	}
	a.DB.Timeout = dbTimeout
	a.DB.MaxConn = dbMaxConn

	a.RPCServer.Port = os.Getenv("RPC_PORT")
	a.RPCServer.Type = os.Getenv("RPC_PROTOCOL")

	a.UserRPC.Port = os.Getenv("USER_RPC_PORT")
	a.UserRPC.Host = os.Getenv("USER_RPC_HOST")
}
