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
	"poptimizer/data/domain/services"
	"poptimizer/data/domain/tables"
	"syscall"
)

type App struct {
	commands chan domain.Command
	events   chan domain.Event
	repo     adapters.Repo
}

func (a *App) initRepo() {
	client := gomoex.NewISSClient(http.DefaultClient)
	factory := tables.NewMainFactory(client)
	a.repo = adapters.Repo{factory}
}

func (a *App) initAdapters() {

}

func (a *App) initCommandBus(ctx context.Context) {
	a.commands = make(chan domain.Command)
}

func (a *App) initEventBus(ctx context.Context) {
	a.events = make(chan domain.Event)
	rules := []domain.Rule{}
	go eventBus(ctx, a.events, a.commands, rules)
}

func (a *App) initDomainServices(ctx context.Context) {
	s := []domain.Service{
		services.WorkStarted{},
	}
	for _, service := range s {
		commands := service.Start(ctx)
		go func() {
			for cmd := range commands {
				a.commands <- cmd
			}
		}()
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

	a.initRepo()
	a.initAdapters()
	a.initCommandBus(ctx)
	a.initEventBus(ctx)
	a.initDomainServices(ctx)

	commandBus(ctx, a.commands, a.events, a.repo)
}
