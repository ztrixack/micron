package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Name     string
	Env      string
	Encoders []EncoderConfig
}

type EncoderConfig struct {
	Encoding string
	Level    string
	Options  struct {
		EncodeLevel string
		Path        string
		FileName    string
		Rotate      bool
		RotateTime  int
		MaxAge      int
	}
}

func (conf *EncoderConfig) LevelEnabler() zap.LevelEnablerFunc {
	highPriority, err := zapcore.ParseLevel(conf.Level)
	if err != nil {
		highPriority = zap.ErrorLevel
	}

	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= highPriority
	})
}

func (conf *EncoderConfig) EncodeLevel() zapcore.LevelEncoder {
	switch conf.Options.EncodeLevel {
	case CapitalColorLevel:
		return zapcore.CapitalColorLevelEncoder
	case CapitalLevel:
		return zapcore.CapitalLevelEncoder
	case LowercaseColorLevel:
		return zapcore.LowercaseColorLevelEncoder
	case LowercaseLevel:
		return zapcore.LowercaseLevelEncoder

	default:
		return zapcore.CapitalLevelEncoder
	}
}

func (conf *EncoderConfig) FilePath() string {
	path := conf.Options.Path
	if len(path) == 0 {
		path = DefaultPath
	}

	if path[:len(path)-1] != "/" {
		path += "/"
	}

	filename := conf.Options.FileName

	if len(filename) == 0 {
		filename = DefaultFilename
	}

	return path + filename
}

func (conf *EncoderConfig) RotateTime() time.Duration {
	rts := conf.Options.RotateTime
	if rts <= 0 {
		rts = DefaultRotateTime
	}

	return time.Duration(rts) * time.Hour
}

func (conf *EncoderConfig) MaxAge() time.Duration {
	rts := conf.Options.MaxAge
	if rts <= 0 {
		rts = DefaultMaxAge
	}

	return time.Duration(rts) * time.Hour
}
