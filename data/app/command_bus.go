package app

import (
	"context"
	"fmt"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
)

func commandBus(
	ctx context.Context,
	commands <-chan domain.Command,
	events chan<- domain.Event,
	repo adapters.Repo,
) {
	for {
		select {
		case cmd := <-commands:
			go handleCommand(ctx, cmd, events, repo)
		case <-ctx.Done():
			return
		}
	}
}

func handleCommand(
	ctx context.Context,
	cmd domain.Command,
	events chan<- domain.Event,
	repo adapters.Repo,
) {
	fmt.Printf("Обработка команды %+v\n", cmd)
	table := repo.Load(cmd.Group, cmd.Name)
	event, err := table.Update(ctx, cmd)
	if err != nil {
		panic("Не удалось обновить таблицу")
	}
	// if event != nil {
	repo.Save(event)
	events <- event
	//}

}
