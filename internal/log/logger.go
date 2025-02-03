package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Logger struct {
	info       *log.Logger
	debug      *log.Logger
	error      *log.Logger
	configName string
}

func New(info, debug, error io.Writer) *Logger {
	l := &Logger{}

	l.SetInfoWriter(info)
	l.SetDebugWriter(debug)
	l.SetErrorWriter(error)

	return l
}

func (l *Logger) SetForConfigFile(fileName string) {
	l.configName = fileName
}

func (l *Logger) SetDebugWriter(debug io.Writer) {
	prefix := fmt.Sprintf("[Debug]: %s: ", l.configName)
	l.debug = log.New(debug, prefix, log.Ltime)
}

func (l *Logger) SetInfoWriter(info io.Writer) {
	prefix := fmt.Sprintf("[Info]: %s: ", l.configName)
	l.info = log.New(info, prefix, log.Ltime)
}

func (l *Logger) SetErrorWriter(error io.Writer) {
	prefix := fmt.Sprintf("[Error]: %s: ", l.configName)
	l.error = log.New(error, prefix, log.Ltime)
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
