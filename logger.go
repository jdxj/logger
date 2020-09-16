package logger

import (
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugar *zap.SugaredLogger
)

func Init(path, mode string) {
	core := zapcore.NewCore(
		encoder(mode),
		syncer(path, mode),
		level(mode),
	)
	sugar = zap.New(core).Sugar()
}

func syncer(path, mode string) zapcore.WriteSyncer {
	switch mode {
	case "debug":
		return zapcore.AddSync(os.Stdout)

	case "release":
		rotation := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 10,
			LocalTime:  false,
			Compress:   false,
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

func Debug(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Info(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warn(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Error(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Sync() {
	_ = sugar.Sync()
}
