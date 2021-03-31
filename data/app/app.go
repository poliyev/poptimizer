package app

import (
	"context"
	"github.com/WLM1ke/gomoex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"sync"
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

	client := gomoex.NewISSClient(http.DefaultClient)
	factory := domain.NewMainFactory(client)
	repo, err := adapters.NewRepo(factory)
	if err != nil {
		panic("не удалось инициализировать репо")
	}

	bus := Bus{repo: repo}

	rules := []domain.Rule{}
	for rule := range rules {
		bus.register(rule)
	}

	services := []domain.Service{
		&domain.WorkStarted{},
	}
	for _, service := range services {
		service.Start(ctx)
		bus.register(service)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		bus.Run(ctx)
	}()

	wg.Wait()
}
