package logs

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type (
	Logger interface {
		Info(...interface{})
		Infof(string, ...interface{})
		Debug(...interface{})
		Debugf(string, ...interface{})
		Error(...interface{})
		Errorf(string, ...interface{})
		Warning(...interface{})
		Warningf(string, ...interface{})
		Fatal(...interface{})
		Fatalf(string, ...interface{})
		Print(...interface{})
		Printf(string, ...interface{})
		Instance() interface{}
	}

	Level     string
	Formatter string

	Option struct {
		Level       Level
		LogFilePath string
		Formatter   Formatter
	}

	logger struct {
		instance *logrus.Logger
	}
)

const (
	Info  Level = "INFO"
	Debug Level = "DEBUG"
	Error Level = "ERROR"

	JSONFormatter Formatter = "JSON"
	TextFormatter Formatter = "TEXT"
)

func (l *logger) Info(args ...interface{}) {
	l.instance.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.instance.Infof(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.instance.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.instance.Debugf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.instance.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.instance.Errorf(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.instance.Warning(args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.instance.Warningf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.instance.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.instance.Fatalf(format, args...)
}

func (l *logger) Print(args ...interface{}) {
	l.instance.Print(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.instance.Printf(format, args...)
}

func (l *logger) Instance() interface{} {
	return l.instance
}

func New(option *Option) (Logger, error) {
	instance := logrus.New()

	if option.Level == Info {
		instance.Level = logrus.InfoLevel
	}

	if option.Level == Debug {
		instance.Level = logrus.DebugLevel
	}

	if option.Level == Error {
		instance.Level = logrus.ErrorLevel
	}

	var formatter logrus.Formatter

	if option.Formatter == JSONFormatter {
		formatter = &logrus.JSONFormatter{}
	} else {
		formatter = &logrus.TextFormatter{}
	}

	instance.Formatter = formatter

	// - check if log file path does exists
	if option.LogFilePath != "" {
		if _, err := os.Stat(option.LogFilePath); os.IsNotExist(err) {
			if _, err = os.Create(option.LogFilePath); err != nil {
				return nil, errors.Wrapf(err, "failed to create log file %s", option.LogFilePath)
			}
		}
		maps := lfshook.PathMap{
			logrus.InfoLevel:  option.LogFilePath,
			logrus.DebugLevel: option.LogFilePath,
			logrus.ErrorLevel: option.LogFilePath,
		}
		instance.Hooks.Add(lfshook.NewHook(maps, formatter))
	}

	return &logger{instance}, nil
}

func DefaultLog() Logger {
	logger, _ := New(&Option{
		Level:       Info,
		Formatter:   TextFormatter,
	})
	return logger
}
