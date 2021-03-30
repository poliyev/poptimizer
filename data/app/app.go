package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"poptimizer/data/services"
	"poptimizer/data/tables"
	"syscall"
)

type App struct {
}

func (a App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	events := make(chan tables.Event)
	commands := make(chan tables.Command)
	InitServices(ctx, commands, events)

	for {
		select {
		case cmd := <-commands:
			go fmt.Printf("%+v\n", cmd)
		case <-ctx.Done():
			return
		}

	}

	//repo, _ := InitRepo()
	//c := time.Tick(time.Minute)
	//fmt.Println("App Started!!!")
	//for t := range c {
	//	table, _ := repo.Get(context.Background(), "trading_dates", "trading_dates")
	//	cmd := tables.Command{}
	//	fmt.Println(t)
	//	fmt.Println(table.Update(context.Background(), cmd))
	//	_ = repo.Save(context.Background(), table)
	//}

}

func InitEventHandlers(commands chan tables.Command) []tables.EventHandler {
	return []tables.EventHandler{
		services.DayStartedHandler{commands},
	}
}

func InitServices(ctx context.Context, commands chan tables.Command, events chan tables.Event) {
	eventHandlers := InitEventHandlers(commands)
	eventBus := EventBus{events, eventHandlers}
	s := []services.Service{
		eventBus,
		services.DayBegins{events},
	}
	for _, service := range s {
		service.Run(ctx)
	}
}
