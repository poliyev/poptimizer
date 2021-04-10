package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
}

func (l Logger) Name() string {
	return "Logger"
}

func (l Logger) Start(ctx context.Context) error {
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

func (l Logger) Shutdown(ctx context.Context) error {
	err := zap.L().Sync()
	if err != nil && errors.Is(err, errors.New("sync /dev/stderr: inappropriate ioctl for device")) {
		return err
	}
	return nil
}
