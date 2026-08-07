package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New(env string) *Logger {
	encoderCfg := buildEncoderConfig(env)

	var encoder zapcore.Encoder
	switch env {
	case "local":
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	case "prod":
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{Logger: zapLogger}
}

func (l *Logger) ErrorRenderPage(err error) {
	l.WithOptions(zap.AddCallerSkip(1)).Error(
		"failed to render page", zap.Error(err),
	)
}

func buildEncoderConfig(env string) zapcore.EncoderConfig {
	var lvlEncoder zapcore.LevelEncoder

	switch env {
	case "local":
		lvlEncoder = zapcore.CapitalColorLevelEncoder
	case "prod":
		lvlEncoder = zapcore.CapitalLevelEncoder
	}

	return zapcore.EncoderConfig{
		NameKey:        "logger",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		TimeKey:        "timestamp",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    lvlEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC1123Z),
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}
