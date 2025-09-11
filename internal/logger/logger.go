package logger

import (
	"log"
)

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
}

type Type int

const (
	Local  Type = iota
	Remote Type = iota
)

var Log Logger

func Init(loggerType Type, source string, hostname string) {
	log.SetFlags(0)
	switch loggerType {
	case Local:
		Log = NewLocalLogger(source, hostname)
	case Remote:
		Log = NewRemoteLogger(source, hostname)
	default:
		Log = NewLocalLogger(source, hostname)
	}
}
