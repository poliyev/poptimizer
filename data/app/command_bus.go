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
			fmt.Printf("Обработка команды %+v\n", cmd)
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
	table, err := repo.Load(ctx, cmd.Group(), cmd.Name())
	if err != nil {
		panic("Не удалось загрузить таблицу")
	}

	table.Update(ctx, cmd)
	for _, event := range table.Events() {
		saveUpdate(ctx, event, repo)
		events <- event
	}
}

func saveUpdate(ctx context.Context, event domain.Event, repo adapters.Repo) {
	if update, ok := event.(domain.TableUpdated); ok {
		err := repo.Save(ctx, update)
		if err != nil {
			panic("Не удалось сохранить таблицу")
		}
	}
}
