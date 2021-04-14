package adapters

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct{}

func NewLogger() *logger {
	return &logger{}
}

func (l logger) Name() string {
	return "logger"
}

func (l logger) Start(ctx context.Context) error {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{
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
	logger, err := config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}

func (l logger) Shutdown(ctx context.Context) error {
	err := zap.L().Sync()
	if err != nil && errors.Is(err, errors.New("sync /dev/stderr: inappropriate ioctl for device")) {
		return err
	}
	return nil
}
