package cfg

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func appTerminationCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 2)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	return ctx
}
