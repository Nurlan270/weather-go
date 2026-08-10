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
	enc := getEncoder(env)

	ws := getWriteSyncer(env)

	enab := getLevelEnabler(env)

	core := zapcore.NewCore(enc, ws, enab)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{Logger: zapLogger}
}

func (l *Logger) ErrorRenderPage(err error) {
	l.WithOptions(zap.AddCallerSkip(1)).Error(
		"failed to render page", zap.Error(err),
	)
}

func getEncoder(env string) zapcore.Encoder {
	encoderCfg := buildEncoderConfig(env)

	var e zapcore.Encoder

	switch env {
	case EnvLocal:
		e = zapcore.NewConsoleEncoder(encoderCfg)
	case EnvProd:
		e = zapcore.NewJSONEncoder(encoderCfg)
	default:
		e = zapcore.NewJSONEncoder(encoderCfg)
	}

	return e
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
	var le zapcore.LevelEnabler

	switch env {
	case EnvLocal:
		le = zapcore.DebugLevel
	case EnvProd:
		le = zapcore.InfoLevel
	default:
		le = zapcore.InfoLevel
	}

	return le
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
