package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Level     zapcore.Level `yaml:"level"`
	Output    string        `yaml:"output"`
	Formatter string        `yaml:"formatter"`
}

var defaultLogger *zap.SugaredLogger

func DefaultLogger() *zap.SugaredLogger {
	return defaultLogger
}

var debugEnabled bool

func DebugEnabled() bool {
	return debugEnabled
}

const timeFormat = "2006-01-02 15:04:05.000"

func init() {
	Init(LogConfig{Level: zap.InfoLevel, Output: "stdout"})
}

func Init(cfg LogConfig) {
	level := zap.NewAtomicLevelAt(cfg.Level)
	debugEnabled = level.Enabled(zap.DebugLevel)
	l, err := zap.Config{
		Level:       level,
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "T",
			LevelKey:      "L",
			NameKey:       "N",
			CallerKey:     "C",
			FunctionKey:   zapcore.OmitKey,
			MessageKey:    "M",
			StacktraceKey: "S",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalColorLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				type appendTimeEncoder interface {
					AppendTimeLayout(time.Time, string)
				}

				if enc, ok := enc.(appendTimeEncoder); ok {
					enc.AppendTimeLayout(t, timeFormat)
					return
				}

				enc.AppendString(t.Format(timeFormat))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{cfg.Output},
		ErrorOutputPaths: []string{cfg.Output},
	}.Build()
	if err != nil {
		panic(err)
	}
	defaultLogger = l.Sugar()
}
