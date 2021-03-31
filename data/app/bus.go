package app

import (
	"context"
	"fmt"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
)

type Bus struct {
	commands  chan domain.Command
	events    chan domain.Event
	consumers []domain.EventConsumer
	repo      *adapters.Repo
}

func (b *Bus) Run(ctx context.Context) {
	b.commands = make(chan domain.Command)
	b.events = make(chan domain.Event)

	for {
		select {
		case cmd := <-b.commands:
			fmt.Printf("Обработка команды %+v\n", cmd)
			go b.handleOneCommand(ctx, cmd)
		case event := <-b.events:
			fmt.Printf("Обработка события %+v\n", event)
			go b.handleOneEvent(ctx, event)
		case <-ctx.Done():
			return
		}
	}
}

func (b *Bus) handleOneCommand(ctx context.Context, cmd domain.Command) {
	table, err := b.repo.Load(ctx, cmd.Group(), cmd.Name())
	if err != nil {
		panic("Не удалось загрузить таблицу")
	}
	for _, event := range table.HandleCommand(ctx, cmd) {
		b.events <- event
	}
}

func (b *Bus) handleOneEvent(ctx context.Context, event domain.Event) {
	if update, ok := event.(domain.TableUpdated); ok {
		err := b.repo.Save(ctx, update)
		if err != nil {
			panic("Не удалось сохранить таблицу")
		}
	}
	for _, consumer := range b.consumers {
		consumer.HandleEvent(ctx, event)
	}
}

func (b *Bus) register(step interface{}) {

	if consumer, ok := step.(domain.EventConsumer); ok {
		b.consumers = append(b.consumers, consumer)
	}

	if cs, ok := step.(domain.CommandSource); ok {
		go func() {
			for cmd := range cs.Commands() {
				b.commands <- cmd
			}
		}()
	}

}
