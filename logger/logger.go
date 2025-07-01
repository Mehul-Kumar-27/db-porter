package logger

import (
	"io"
	"log"
	"os"
	"strings"
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

// formatLabel formats and ensures the label looks like [LABEL]
func formatLabel(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return ""
	}
	return "[" + label + "] "
}

// New creates a new Logger instance with a custom label
func New(label string, out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}

	prefix := formatLabel(label)

	return &Logger{
		infoLog:    log.New(out, colorBlue+prefix+"INFO: "+colorReset, log.LstdFlags|log.Lshortfile),
		warningLog: log.New(out, colorYellow+prefix+"WARNING: "+colorReset, log.LstdFlags|log.Lshortfile),
		errorLog:   log.New(out, colorRed+prefix+"ERROR: "+colorReset, log.LstdFlags|log.Lshortfile),
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

// Optional default logger
var defaultLogger = New("DEFAULT", os.Stdout)

func Info(format string, v ...any) {
	defaultLogger.Info(format, v...)
}

func Warning(format string, v ...any) {
	defaultLogger.Warning(format, v...)
}

func Error(format string, v ...any) {
	defaultLogger.Error(format, v...)
}
