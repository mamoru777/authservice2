package main

import (
	"github.com/caarlos0/env/v8"
	"github.com/mamoru777/authservice2/internal/app"
	"github.com/mamoru777/authservice2/internal/config"
	"github.com/mamoru777/authservice2/internal/mylogger"
	"github.com/sirupsen/logrus"
)

func main() {

	logger := mylogger.New(&logrus.Logger{})
	file := mylogger.ILogger.Init(logger) //mylogger.Init()
	defer file.Close()
	dataBaseconfig := config.DataBaseConfig{}
	grpcServerConfig := config.GrpcServerConfig{}
	if err := env.Parse(&dataBaseconfig); err != nil {
		logger.Logger.Fatalf("ошибка при получении переменных окружения, %v", err)
	}
	if err := env.Parse(&grpcServerConfig); err != nil {
		logger.Logger.Fatalf("ошибка при получении переменных окружения, %v", err)
	}
	if err := app.Run(dataBaseconfig, grpcServerConfig, logger); err != nil {
		logger.Logger.Fatal("ошибка при запуске сервера ", err)
	}
}
