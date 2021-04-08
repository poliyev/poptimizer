package app

import (
	"context"
	"fmt"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"sync"
)

type Bus struct {
	commands  chan domain.Command
	events    chan domain.Event
	consumers []chan domain.Event
	repo      *adapters.Repo
	ctx       context.Context
	wg        sync.WaitGroup
}

func (b *Bus) Run(ctx context.Context) {
	if b.ctx != nil {
		panic("Шина уже запущена")
	}

	b.events = make(chan domain.Event)
	b.ctx = ctx

	for {
		b.wg.Add(1)
		select {
		case cmd := <-b.commands:
			fmt.Printf("Обработка команды %+v\n", cmd)
			go func() {
				defer b.wg.Done()
				b.handleOneCommand(ctx, cmd)
			}()
		case event := <-b.events:
			fmt.Printf("Обработка события %+v\n", event)
			go func() {
				defer b.wg.Done()
				go b.handleOneEvent(ctx, event)
			}()
		case <-ctx.Done():
			b.wg.Done()
			b.wg.Wait()
			return
		}
	}
}

func (b *Bus) handleOneCommand(ctx context.Context, cmd domain.Command) {
	table, err := b.repo.Load(ctx, cmd.ID())
	if err != nil {
		panic("Не удалось загрузить таблицу")
	}
	for _, event := range table.HandleCommand(ctx, cmd) {
		b.events <- event
	}
}

func (b *Bus) handleOneEvent(ctx context.Context, event domain.Event) {
	err := b.repo.Save(ctx, event)
	if err != nil {
		panic("Не удалось сохранить таблицу")
	}
	for _, consumer := range b.consumers {
		consumer <- event
	}
}

func (b *Bus) register(step interface{}) {
	if b.commands == nil {
		b.commands = make(chan domain.Command)
	}
	// Добавить логинг о старте и завершении
	if consumer, ok := step.(domain.EventConsumer); ok {
		newChan := make(chan domain.Event)
		b.consumers = append(b.consumers, newChan)

		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			consumer.StartHandleEvent(b.ctx, newChan)
		}()
	}
	// Добавить логинг о старте и завершении
	if source, ok := step.(domain.CommandSource); ok {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			source.StartProduceCommands(b.ctx, b.commands)
		}()
	}

}
