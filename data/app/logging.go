package app

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func StartLogging() {
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
		zap.L().Panic("Не удалось запустить логирование")
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("Логирование запущено")
}

func ShutdownLogging() {
	zap.L().Info("Логирование остановлено")
	err := zap.L().Sync()
	if err != nil && errors.Is(err, errors.New("sync /dev/stderr: inappropriate ioctl for device")) {
		zap.L().Error("Не удалось остановить логирование", zap.Error(err))
	}
}
