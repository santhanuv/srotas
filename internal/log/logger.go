package log

import (
	"bytes"
	"encoding/json"
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
	l := &Logger{}

	l.SetInfoWriter(info)
	l.SetDebugWriter(debug)
	l.SetErrorWriter(error)

	return l
}

func (l *Logger) SetDebugWriter(debug io.Writer) {
	l.debug = log.New(debug, "[Debug]: ", log.Ltime)
}

func (l *Logger) SetInfoWriter(info io.Writer) {
	l.info = log.New(info, "[Info]: ", log.Ltime)
}

func (l *Logger) SetErrorWriter(error io.Writer) {
	l.error = log.New(error, "[Error]: ", log.Ltime)
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

func (l *Logger) DebugJson(v []byte, prefix string, args ...any) {
	var buf bytes.Buffer
	json.Indent(&buf, v, "", " ")

	l.Debug("%s\n%s", prefix, buf.String())
}
