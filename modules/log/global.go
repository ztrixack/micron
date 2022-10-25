package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	FileEncoder        = "file"
	ConsoleEncoder     = "console"
	StackdriverEncoder = "stackdriver"

	DefaultPath         = "log/"
	DefaultFilename     = "app-%Y-%m-%d.log"
	DefaultRotateTime   = 24
	DefaultMaxAge       = 90 * 24
	CapitalColorLevel   = "capital-color"
	CapitalLevel        = "capital"
	LowercaseColorLevel = "lowercase-color"
	LowercaseLevel      = "lowercase"
)

var logger *zap.Logger

func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Critical(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
	os.Exit(1)
}
