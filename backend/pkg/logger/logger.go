package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.infoLogger.Println(append([]interface{}{msg}, keysAndValues...)...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.errorLogger.Println(append([]interface{}{msg}, keysAndValues...)...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.errorLogger.Fatal(append([]interface{}{msg}, keysAndValues...)...)
}
