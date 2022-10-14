package log

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Init(conf Config) {
	if conf.Development {
		logger = development(conf)
	} else {
		logger = production(conf)
	}

	defer logger.Sync()

}

func production(conf Config) *zap.Logger {
	if conf.File {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoder := zapcore.NewJSONEncoder(config)
		consoleEncoder := zapcore.NewConsoleEncoder(config)
		logFile, _ := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer := zapcore.AddSync(logFile)
		defaultLogLevel := zapcore.DebugLevel
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)

		return zap.New(core, zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		logger, err := zap.NewProduction(zap.AddCallerSkip(1))
		if err != nil {
			log.Fatalf(err.Error())
		}

		return logger
	}

}

func development(conf Config) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.DisableStacktrace = true
	logger, _ := config.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))

	return logger
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Event(msg string, event string, fields ...zapcore.Field) {
	fields = append([]zapcore.Field{zap.String("event", event)}, fields...)
	logger.Info(msg, fields...)
}

func Warn(msg string, reason string, fields ...zapcore.Field) {
	fields = append([]zapcore.Field{zap.String("reason", reason)}, fields...)
	logger.Warn(msg, fields...)
}

func Error(msg string, err error, fields ...zapcore.Field) {
	fields = append(fields, zap.Error(err))
	logger.Error(msg, fields...)
	os.Exit(1)
}
