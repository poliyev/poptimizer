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
	commands chan domain.Command
	events   chan domain.Event
	repo     adapters.Repo
}

func (a *App) initRepo() {
	client := gomoex.NewISSClient(http.DefaultClient)
	factory := domain.NewMainFactory(client)
	repo, err := adapters.NewRepo(factory)
	if err != nil {
		panic("не удалось инициализировать репо")
	}
	a.repo = *repo
}

func (a *App) initAdapters() {

}

func (a *App) initDomainServices(ctx context.Context) {
	s := []domain.Service{
		&domain.WorkStarted{},
	}
	for _, service := range s {
		service.Start(ctx)

		if source, ok := service.(domain.CommandSource); ok {
			go func() {
				for cmd := range source.Commands() {
					a.commands <- cmd
				}
			}()
		}
	}
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

	a.commands = make(chan domain.Command)
	a.events = make(chan domain.Event)
	rules := []domain.Rule{}

	a.initRepo()
	a.initAdapters()
	a.initDomainServices(ctx)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		commandBus(ctx, a.commands, a.events, a.repo)
	}()

	go func() {
		defer wg.Done()
		go eventBus(ctx, a.events, a.commands, rules)
	}()

	wg.Wait()
}
