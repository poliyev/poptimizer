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
			newCmd := processRules(event, rules)
			sendCommands(newCmd, commands)
		case <-ctx.Done():
			return
		}
	}
}

func processRules(event domain.Event, rules []domain.Rule) []domain.Command {
	cmd := make([]domain.Command, 0)
	for _, rule := range rules {
		if rule.Match(event) {
			cmd = append(cmd, rule.Handle(event)...)
		}
	}

	return cmd
}

func sendCommands(newCmd []domain.Command, commands chan<- domain.Command) {
	for _, cmd := range newCmd {
		commands <- cmd
	}
}
