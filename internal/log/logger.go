package log

import (
	"fmt"
	"io"
	"log"
)

type Logger struct {
	info  *log.Logger
	debug *log.Logger
	error *log.Logger
}

func New(info, debug, error io.Writer) *Logger {
	infoLogger := log.New(info, "[Info]: ", log.Ltime)
	debugLogger := log.New(debug, "[Debug]: ", log.Ltime)
	errorLogger := log.New(error, "[Error]: ", log.Ltime)

	return &Logger{
		info:  infoLogger,
		debug: debugLogger,
		error: errorLogger,
	}
}

func (l *Logger) Info(format string, args ...any) {
	formatLine := fmt.Sprintf("%s\n", format)
	l.info.Printf(formatLine, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	formatLine := fmt.Sprintf("%s\n", format)
	l.debug.Printf(formatLine, args...)
}

func (l *Logger) Error(format string, args ...any) {
	formatLine := fmt.Sprintf("%s\n", format)
	l.error.Printf(formatLine, args...)
}

func (l *Logger) Fatal(format string, args ...any) {
	formatLine := fmt.Sprintf("%s\n", format)
	l.error.Fatalf(formatLine, args...)
}
