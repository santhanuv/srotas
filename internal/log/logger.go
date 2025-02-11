package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Logger is a structured logging utility for the CLI application, providing
// categorized logging for informational, debug, and error messages.
// It maintains a configuration name for context and a debugMode flag to
// enable or disable debug logging.
type Logger struct {
	info       *log.Logger // Logger for information type logs.
	debug      *log.Logger // Logger for debug type logs.
	error      *log.Logger // Logger for error and fatal type logs.
	configName string      // Name of the config for contextual logging.
	debugMode  bool        // Indicates whether debug logging is enabled.
}

// New creates a new [Logger].
// The info, debug, error variable defines the output destination for the info, debug, and error logs, respectively.
// By default the debugMode will be set to false.
func New(info, debug, error io.Writer) *Logger {
	l := &Logger{}

	l.configName = "config"
	l.debugMode = false
	l.info = log.New(info, "[INFO]: ", log.Ltime)
	l.debug = log.New(debug, "[DEBUG]: ", log.Ltime)
	l.error = log.New(error, "[Error]: ", log.Ltime)

	return l
}

// SetConfig set the config file for the logger.
func (l *Logger) SetConfig(fileName string) {
	l.configName = fileName
}

// SetDebugOutput set the output destination for debug logs.
func (l *Logger) SetDebugOutput(debug io.Writer) {
	l.debug.SetOutput(debug)
}

// SetInfoOutput set the output destination for info logs.
func (l *Logger) SetInfoOutput(info io.Writer) {
	l.info.SetOutput(info)
}

// SetErrorOutput set the output destination for error logs.
func (l *Logger) SetErrorOutput(error io.Writer) {
	l.error.SetOutput(error)
}

// SetDebugMode enables or disables debug logging based on the provided value.
// When set to true, debug logs are enabled; when set to false, they are suppressed.
func (l *Logger) SetDebugMode(debugMode bool) {
	l.debugMode = debugMode
}

// Info logs the given message as an info log.
// Arguments are handled in the manner of Printf.
func (l *Logger) Info(format string, args ...any) {
	formatLine := fmt.Sprintf("%s: %s\n", l.configName, format)
	l.info.Printf(formatLine, args...)
}

// Debug logs the given message as a debug log.
// Arguments are handled in the manner of Printf.
func (l *Logger) Debug(format string, args ...any) {
	if !l.debugMode {
		return
	}

	formatLine := fmt.Sprintf("%s: %s\n", l.configName, format)
	l.debug.Printf(formatLine, args...)
}

// Error logs the given message as an error log.
// Arguments are handled in the manner of Printf.
func (l *Logger) Error(format string, args ...any) {
	formatLine := fmt.Sprintf("%s: %s\n", l.configName, format)
	l.error.Printf(formatLine, args...)
}

// Fatal logs the given message as an error log and exits the program.
// Arguments are handled in the manner of Printf.
func (l *Logger) Fatal(format string, args ...any) {
	formatLine := fmt.Sprintf("%s: %s\n", l.configName, format)
	l.error.Fatalf(formatLine, args...)
}

// DebugJson logs the given JSON value with a formatted prefix at the start of the line as debug log.
// The formatted prefix is constructed using the provided prefix string and args, which are handled in the same manner as Printf.
func (l *Logger) DebugJson(v []byte, prefix string, args ...any) {
	if !l.debugMode {
		return
	}

	var buf bytes.Buffer
	json.Indent(&buf, v, "", " ")

	l.Debug("%s\n%s", prefix, buf.String())
}
