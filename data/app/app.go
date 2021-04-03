package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"syscall"
)

type App struct {
}

func (a *App) initAdapters() {

}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Printf("\n[WARN] interrupt signal")
		cancel()
	}()

	iss := adapters.NewISSClient()
	factory := domain.NewMainFactory(iss)
	repo := adapters.NewRepo(factory)

	bus := Bus{repo: repo}

	rules := []domain.Rule{}
	for _, rule := range rules {
		bus.register(rule)
	}

	sources := []domain.CommandSource{
		&domain.CheckTradingDay{},
	}
	for _, source := range sources {
		bus.register(source)
	}

	bus.Run(ctx)
}
