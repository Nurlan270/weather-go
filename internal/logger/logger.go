package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Logger struct {
	*zap.Logger
}

func New(env string) *Logger {
	encoder := getEncoder(env)

	ws := getWriteSyncer(env)

	enabler := getLevelEnabler(env)

	core := zapcore.NewCore(encoder, ws, enabler)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{Logger: zapLogger}
}

func (l *Logger) ErrorRenderPage(err error) {
	l.WithOptions(zap.AddCallerSkip(1)).Error(
		"failed to render page", zap.Error(err),
	)
}

// Printf is used by BigCache package.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

func getEncoder(env string) zapcore.Encoder {
	encoderCfg := buildEncoderConfig(env)

	var encoder zapcore.Encoder

	switch env {
	case EnvLocal:
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	case EnvProd:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	return encoder
}

func getWriteSyncer(env string) zapcore.WriteSyncer {
	var ws zapcore.WriteSyncer

	switch env {
	case EnvLocal:
		ws = zapcore.AddSync(os.Stdout)
	case EnvProd:
		ws = zapcore.Lock(buildFileWriteSyncer())
	default:
		ws = zapcore.AddSync(os.Stdout)
	}

	return ws
}

func getLevelEnabler(env string) zapcore.LevelEnabler {
	var enabler zapcore.LevelEnabler

	switch env {
	case EnvLocal:
		enabler = zapcore.DebugLevel
	case EnvProd:
		enabler = zapcore.InfoLevel
	default:
		enabler = zapcore.InfoLevel
	}

	return enabler
}

func buildEncoderConfig(env string) zapcore.EncoderConfig {
	var levelEncoder zapcore.LevelEncoder

	switch env {
	case EnvLocal:
		levelEncoder = zapcore.CapitalColorLevelEncoder
	case EnvProd:
		levelEncoder = zapcore.CapitalLevelEncoder
	default:
		levelEncoder = zapcore.CapitalLevelEncoder
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
		EncodeLevel:    levelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}

func buildFileWriteSyncer() zapcore.WriteSyncer {
	logPath := filepath.Join(
		"logs", fmt.Sprintf("%s.log", time.Now().Format(time.DateOnly)),
	)

	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	return zapcore.AddSync(file)
}
