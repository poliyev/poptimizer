package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger модуль логирования на основе глобального логера zap.
type Logger struct{}

// NewLogger создает новый модуль логирования.
func NewLogger() *Logger {
	return &Logger{}
}

// Start устанавливает настройки глобального логера zap.
func (l Logger) Start(_ context.Context) error {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "M",
		LevelKey:      "L",
		TimeKey:       "T",
		NameKey:       "N",
		CallerKey:     "C",
		FunctionKey:   zapcore.OmitKey,
		StacktraceKey: "S",
		LineEnding:    "\n",

		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       zapcore.RFC3339TimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("logger start failed: %w", err)
	}

	zap.ReplaceGlobals(logger)

	return nil
}

var errStderrSync = errors.New("sync /dev/stderr: inappropriate ioctl for device")

// Shutdown синхронизирует записи в лог.
func (l Logger) Shutdown(_ context.Context) error {
	err := zap.L().Sync()
	if err != nil && errors.Is(err, errStderrSync) {
		return fmt.Errorf("logger shutdown error: %w", err)
	}

	return nil
}

// TypeField - поле для логера с коротким типом объекта (с убранным знаком указателя и путем к нему).
func TypeField(value interface{}) zap.Field {
	parts := strings.Split(fmt.Sprintf("%T", value), ".")
	return zap.String("type", parts[len(parts)-1])
}
