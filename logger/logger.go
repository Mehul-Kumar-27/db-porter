package logger

import (
	"io"
	"log"
	"os"
)

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Logger struct holds the log configurations
type Logger struct {
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
}

// New creates a new Logger instance
func New(out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout // Default to stdout
	}

	return &Logger{
		infoLog:    log.New(out, colorBlue+"INFO: "+colorReset, log.LstdFlags|log.Lshortfile),
		warningLog: log.New(out, colorYellow+"WARNING: "+colorReset, log.LstdFlags|log.Lshortfile),
		errorLog:   log.New(out, colorRed+"ERROR: "+colorReset, log.LstdFlags|log.Lshortfile),
	}
}

// Info logs informational messages (blue)
func (l *Logger) Info(format string, v ...any) {
	l.infoLog.Printf(format, v...)
}

// Warning logs warning messages (yellow)
func (l *Logger) Warning(format string, v ...any) {
	l.warningLog.Printf(format, v...)
}

// Error logs error messages (red)
func (l *Logger) Error(format string, v ...any) {
	l.errorLog.Printf(format, v...)
}

// Default logger instance (optional)
var defaultLogger = New(os.Stdout)

// Package-level convenience functions (optional)
func Info(format string, v ...any) {
	defaultLogger.Info(format, v...)
}

func Warning(format string, v ...any) {
	defaultLogger.Warning(format, v...)
}

func Error(format string, v ...any) {
	defaultLogger.Error(format, v...)
}
