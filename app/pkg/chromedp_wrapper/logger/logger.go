package logger

import (
	"log"
	"os"
)

type logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger() *logger {
	flags := log.LstdFlags | log.Lshortfile
	infoLogger := log.New(os.Stdout, "INFO: ", flags)
	warnLogger := log.New(os.Stdout, "WARN: ", flags)
	errorLogger := log.New(os.Stdout, "ERROR: ", flags)
	return &logger{infoLogger: infoLogger, warnLogger: warnLogger, errorLogger: errorLogger}
}

func (l *logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.warnLogger.Println(v...)
}

func (l *logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}
