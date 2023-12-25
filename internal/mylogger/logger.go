package mylogger

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/mediasoft-internship/internship/mamoru777/foundation/loginit"
	"os"
)

type Logger struct {
	Logger *logrus.Logger
}

func New(logger *logrus.Logger) *Logger {
	return &Logger{
		Logger: logger,
	}
}

func (l *Logger) Init() *os.File {
	l.Logger = loginit.LogInit("tcp", "localhost:5044", "AuthService")
	_ = l.Logger
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		l.Logger.Out = file
	} else {
		logrus.Error("Не удалось логировать в файл, использую логирование в консоль")
	}
	return file
}
