package app

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func TerminationSignal() <-chan struct{} {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 2)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		zap.L().Info("Поступил сигнал остановки")
		cancel()
	}()

	return ctx.Done()
}
