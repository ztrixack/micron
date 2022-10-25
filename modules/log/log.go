package log

import (
	"context"
	"io"
	"os"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/ztrixack/micron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogContext interface {
	Build()
}

type logContextImpl struct {
	config Config
	cores  []zapcore.Core
}

func New(conf Config) LogContext {
	cores := []zapcore.Core{}

	for _, encoderConf := range conf.Encoders {
		switch encoderConf.Encoding {
		case FileEncoder:
			writer := getFileWriter(encoderConf)
			encoder := zap.NewProductionEncoderConfig()
			fileEncoder := zapcore.NewJSONEncoder(encoder)
			core := zapcore.NewCore(fileEncoder, zapcore.AddSync(writer), encoderConf.LevelEnabler())
			cores = append(cores, core)

		case ConsoleEncoder:
			encoder := zap.NewDevelopmentEncoderConfig()
			encoder.EncodeLevel = encoderConf.EncodeLevel()
			consoleEncoder := zapcore.NewConsoleEncoder(encoder)
			core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), encoderConf.LevelEnabler())
			cores = append(cores, core)

		case StackdriverEncoder:
			encoder := zapcore.EncoderConfig{
				LevelKey:       "severity",
				NameKey:        "logger",
				CallerKey:      "caller",
				StacktraceKey:  "stacktrace",
				TimeKey:        "time",
				MessageKey:     "message",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
				EncodeLevel:    encodeLevel,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}
			stackdriverEncoder := zapcore.NewJSONEncoder(encoder)
			core := zapcore.NewCore(stackdriverEncoder, zapcore.AddSync(os.Stdout), encoderConf.LevelEnabler())
			cores = append(cores, core)
		}
	}

	return &logContextImpl{
		config: conf,
		cores:  cores,
	}
}

func getFileWriter(encoderConf EncoderConfig) io.Writer {
	if encoderConf.Options.Rotate {
		rotator, err := rotatelogs.New(
			encoderConf.FilePath(),
			rotatelogs.WithMaxAge(encoderConf.MaxAge()),
			rotatelogs.WithRotationTime(encoderConf.RotateTime()))
		micron.AppCtx.AddTerminateFunc(func(ctx context.Context) {
			if err := rotator.Close(); err != nil {
				panic(err)
			}
		})
		if err != nil {
			panic(err)
		}

		return rotator
	} else {
		file, err := os.OpenFile(encoderConf.FilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		micron.AppCtx.AddTerminateFunc(func(ctx context.Context) {
			if err := file.Close(); err != nil {
				panic(err)
			}
		})
		if err != nil {
			panic(err)
		}

		return file
	}
}

func (ctx *logContextImpl) AddZapcore(core zapcore.Core) {
	ctx.cores = append(ctx.cores, core)
}

func (ctx *logContextImpl) Build() {
	logcore := zapcore.NewTee(ctx.cores...)
	logger = zap.New(logcore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("Debug") // logging.Debug.String()
	case zapcore.InfoLevel:
		enc.AppendString("Info") // logging.Info.String()
	case zapcore.WarnLevel:
		enc.AppendString("Warning") // logging.Warning.String()
	case zapcore.ErrorLevel:
		enc.AppendString("Error") // logging.Error.String()
	case zapcore.DPanicLevel:
		enc.AppendString("Critical") // logging.Critical.String()
	case zapcore.PanicLevel:
		enc.AppendString("Alert") // logging.Alert.String()
	case zapcore.FatalLevel:
		enc.AppendString("Emergency") // logging.Emergency.String()
	}
}
