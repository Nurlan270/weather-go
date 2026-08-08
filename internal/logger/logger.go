package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"time"
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
	case "local":
		e = zapcore.NewConsoleEncoder(encoderCfg)
	case "prod":
		e = zapcore.NewJSONEncoder(encoderCfg)
	default:
		e = zapcore.NewJSONEncoder(encoderCfg)
	}

	return e
}

func getWriteSyncer(env string) zapcore.WriteSyncer {
	var ws zapcore.WriteSyncer
	switch env {
	case "local":
		ws = zapcore.AddSync(os.Stdout)
	case "prod":
		ws = zapcore.Lock(buildFileWriteSyncer())
	default:
		ws = zapcore.AddSync(os.Stdout)
	}

	return ws
}

func getLevelEnabler(env string) zapcore.LevelEnabler {
	var le zapcore.LevelEnabler
	switch env {
	case "local":
		le = zapcore.DebugLevel
	case "prod":
		le = zapcore.InfoLevel
	default:
		le = zapcore.InfoLevel
	}

	return le
}

func buildEncoderConfig(env string) zapcore.EncoderConfig {
	var lvlEncoder zapcore.LevelEncoder

	switch env {
	case "local":
		lvlEncoder = zapcore.CapitalColorLevelEncoder
	case "prod":
		lvlEncoder = zapcore.CapitalLevelEncoder
	default:
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
