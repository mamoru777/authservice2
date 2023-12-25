package mylogger

import "os"

//go:generate mockgen -destination mock_logger.go -package mylogger . ILogger

type ILogger interface {
	Init() *os.File
}
