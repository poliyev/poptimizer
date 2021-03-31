package app

import (
	"context"
	"fmt"
	"poptimizer/data/domain"
)

func eventBus(ctx context.Context, events <-chan domain.Event, commands chan<- domain.Command, rules []domain.Rule) {
	for {
		select {
		case event := <-events:
			fmt.Printf("Обработка события %+v\n", event)
			processRules(event, rules, commands)
		case <-ctx.Done():
			return
		}
	}
}

func processRules(event domain.Event, rules []domain.Rule, commands chan<- domain.Command) []domain.Command {
	cmd := make([]domain.Command, 0)
	for _, rule := range rules {
		rule.HandleEvent(context.TODO(), event)
		sendCommands(rule, commands)
	}

	return cmd
}

func sendCommands(rule domain.Rule, commands chan<- domain.Command) {
	for cmd := range rule.Commands() {
		commands <- cmd
	}
}
