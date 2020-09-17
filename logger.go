package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	desugar *zap.Logger
	sugar   *zap.SugaredLogger
)

type OptionFunc func(opts *Options)

type Options struct {
	Mode     string
	FileName string

	MaxSize    int
	MaxAge     int
	MaxBackups int

	LocalTime bool
	Compress  bool
}

func NewPathMode(path, mode string) {
	optF := OptionFunc(func(opts *Options) {
		opts.Mode = mode
		opts.FileName = path
		opts.MaxSize = 500
		opts.MaxAge = 30
		opts.MaxBackups = 10
		opts.LocalTime = false
		opts.Compress = false
	})
	New(optF)
}

func New(optsF ...OptionFunc) {
	opts := new(Options)
	for _, optF := range optsF {
		optF(opts)
	}

	core := zapcore.NewCore(
		encoder(opts.Mode),
		syncer(opts),
		level(opts.Mode),
	)
	desugar = zap.New(core)
	sugar = desugar.Sugar()
}

func syncer(opts *Options) zapcore.WriteSyncer {
	switch opts.Mode {
	case "debug":
		return zapcore.AddSync(os.Stdout)

	case "release":
		rotation := &lumberjack.Logger{
			Filename:   opts.FileName,
			MaxSize:    opts.MaxSize,
			MaxAge:     opts.MaxAge,
			MaxBackups: opts.MaxBackups,
			LocalTime:  opts.LocalTime,
			Compress:   opts.Compress,
		}
		return zapcore.AddSync(rotation)
	}
	return nil
}

func encoder(mode string) zapcore.Encoder {
	switch mode {
	case "debug":
		return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	case "release":
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
	return nil
}

func level(mode string) zapcore.Level {
	switch mode {
	case "release":
		return zap.InfoLevel
	}
	return zap.DebugLevel
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Sync() {
	_ = desugar.Sync()
}
