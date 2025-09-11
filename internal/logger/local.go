package logger

import (
	"log"
	"time"
)

type LocalLogger struct {
	Source   string
	Hostname string
}

func NewLocalLogger(source string, hostname string) *LocalLogger {
	return &LocalLogger{
		Source:   source,
		Hostname: hostname,
	}
}

func (ll *LocalLogger) LogWithLevel(message string, level string) {
	timeNow := time.Now().Format("15:04:05")
	log.Printf("[%s] [%s] [%s] [%s] %s\n", timeNow, level, ll.Source, ll.Hostname, message)
}

func (ll *LocalLogger) Debug(message string) {
	ll.LogWithLevel(message, "DEBUG")
}

func (ll *LocalLogger) Info(message string) {
	ll.LogWithLevel(message, "INFO")
}

func (ll *LocalLogger) Warn(message string) {
	ll.LogWithLevel(message, "WARN")
}

func (ll *LocalLogger) Error(msg string) {
	ll.LogWithLevel(msg, "ERROR")
}
