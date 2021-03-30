package app

import (
	"context"
	"fmt"
	"poptimizer/data/tables"
)

type EventBus struct {
	events   chan tables.Event
	handlers []tables.EventHandler
}

func (e EventBus) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case event := <-e.events:
				fmt.Printf("%+v\n", event)
				for _, handler := range e.handlers {
					if handler.Match(event) {
						handler.Handle(event)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
