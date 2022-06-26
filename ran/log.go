package ran

import (
	"fmt"
	"io"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Error
	Discard
)

func NewLogLevel(level string) (LogLevel, error) {
	l, ok := map[string]LogLevel{
		"debug":   Debug,
		"info":    Info,
		"error":   Error,
		"discard": Discard,
	}[level]
	if !ok {
		return 0, fmt.Errorf("unknow log level: %q", level)
	}
	return l, nil
}

func (l LogLevel) String() string {
	return map[LogLevel]string{
		Debug:   "DEBUG",
		Info:    "INFO",
		Error:   "ERROR",
		Discard: "DISCARD",
	}[l]
}

type StdLogger struct {
	w     io.Writer
	level LogLevel
}

func NewStdLogger(w io.Writer, level LogLevel) Logger {
	return &StdLogger{w, level}
}

func (s *StdLogger) Debug(msg string, args ...any) {
	s.print(Debug, msg, args...)
}

func (s *StdLogger) Info(msg string, args ...any) {
	s.print(Info, msg, args...)

}

func (s *StdLogger) Error(msg string, args ...any) {
	s.print(Error, msg, args...)

}

func (s *StdLogger) print(level LogLevel, msg string, args ...any) {
	if level < s.level {
		return
	}

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	msg = fmt.Sprintf("%s", msg)

	if len(args) == 0 {
		fmt.Fprint(s.w, msg)
		return
	}
	fmt.Fprintf(s.w, msg, args...)
}
