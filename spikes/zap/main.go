package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Counter()
	trace   = kingpin.Flag("trace", "Enable trace logging.").Bool()
)

func main() {
	kingpin.Parse()

	var level zapcore.Level
	switch {
	case *verbose >= 3:
		level = zapcore.DebugLevel
	case *verbose == 2:
		level = zapcore.InfoLevel
	case *verbose == 1:
		level = zapcore.WarnLevel
	default:
		level = zapcore.ErrorLevel
	}

	var logger *zap.SugaredLogger
	if *trace {
		tmpfile, err := ioutil.TempFile("", "zap-demo")
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to open trace log:", err)
		}
		fmt.Fprintln(os.Stderr, "logging to:", tmpfile.Name())
		logger = NewLogger(level, tmpfile).Sugar()
	} else {
		logger = NewLogger(level, nil).Sugar()
	}
	defer logger.Sync()

	logger.Errorw("an error",
		"metavar", "foo")
	logger.Warnw("a warning",
		"metavar", "bar",
		"question", "6 * 9")
	logger.With(
		"metavar", "baz",
		"answer", 42).
		Infow("some info")
	logger.With(
		"hint", "The secret is to bang the rocks together, guys.").
		Debugf("debug=%d", *verbose)
}

func NewLogger(level zapcore.Level, traceLog zapcore.WriteSyncer) *zap.Logger {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:     "msg",
			NameKey:        "logger",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}),
		os.Stderr,
		level,
	)
	if traceLog != nil {
		trace := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				MessageKey:     "msg",
				LevelKey:       "level",
				NameKey:        "logger",
				TimeKey:        "time",
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
			}),
			traceLog,
			zapcore.DebugLevel,
		)
		core = zapcore.NewTee(core, trace)
	}
	return zap.New(core)
}
