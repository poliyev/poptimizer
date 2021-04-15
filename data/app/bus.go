package app

import (
	"context"
	"fmt"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Bus - шина для обработки команд и событий, поддерживающую интерфейс модуля приложения.
//
// Регистрирует все шаги бизнес-логики — источники команд, правила обработки доменных событий и потребителей событий.
// Процесс работы шины поддерживает данные в актуальном состоянии.
type Bus struct {
	repo             adapters.TableRepo
	handlersTimeouts time.Duration

	commands  chan domain.Command
	events    chan domain.Event
	consumers []chan domain.Event

	wg          sync.WaitGroup
	loopCtx     context.Context
	loopCancel  context.CancelFunc
	loopStopped chan struct{}
}

// NewBus - создает шину для обработки команд и событий.
func NewBus(repo adapters.TableRepo, eventBusTimeouts time.Duration) *Bus {
	ctx, cancel := context.WithCancel(context.Background())

	bus := Bus{
		repo:             repo,
		handlersTimeouts: eventBusTimeouts,

		commands: make(chan domain.Command),
		events:   make(chan domain.Event),

		loopCtx:     ctx,
		loopCancel:  cancel,
		loopStopped: make(chan struct{}),
	}
	steps := []interface{}{
		// Источники команд
		domain.NewCheckTradingDay(),
		// Правила

		// Потребители сообщений
	}

	for _, step := range steps {
		bus.register(step)
	}

	return &bus
}

// Name - модуль приложения Bus.
func (b *Bus) Name() string {
	return "Bus"
}

// Start - запускает основной цикл обработки команд и событий.
func (b *Bus) Start(_ context.Context) error {
	go func() {
		b.loop(b.loopCtx)
	}()

	return nil
}

// Shutdown - завершает работу основного цикла обработки команд и событий.
func (b *Bus) Shutdown(ctx context.Context) error {
	b.loopCancel()

	select {
	case <-b.loopStopped:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("bus shutdown error: %w", ctx.Err())
	}
}

// loop осуществляет вызов обработчиков команд и событий при их наличии до отмены контекста цикла.
func (b *Bus) loop(ctx context.Context) {
	for {
		b.wg.Add(1)
		select {
		case cmd := <-b.commands:
			go func() {
				defer b.wg.Done()
				b.handleOneCommand(ctx, cmd)
			}()
		case event := <-b.events:
			go func() {
				defer b.wg.Done()

				go b.handleOneEvent(ctx, event)
			}()
		case <-ctx.Done():
			b.wg.Done()
			b.wg.Wait()
			close(b.loopStopped)

			return
		}
	}
}

// handleOneCommand - загружает таблицу, вызывает обработчик команды и направляет возникшие события в очередь.
func (b *Bus) handleOneCommand(ctx context.Context, cmd domain.Command) {
	zap.L().Info("Command", zap.Stringer("table", cmd.ID()))

	ctx, cancel := context.WithTimeout(ctx, b.handlersTimeouts)

	defer cancel()

	table, err := b.repo.Load(ctx, cmd.ID())
	if err != nil {
		zap.L().Panic("Command", zap.Stringer("load", cmd.ID()), zap.Error(err))
	}

	for _, event := range table.HandleCommand(ctx, cmd) {
		b.events <- event
	}
}

// handleOneEvent - сохраняет событие, а после этого рассылает его всем потребителям событий.
func (b *Bus) handleOneEvent(ctx context.Context, event domain.Event) {
	zap.L().Info("Event", zap.Stringer("table", event.ID()))

	ctx, cancel := context.WithTimeout(ctx, b.handlersTimeouts)

	defer cancel()

	if err := b.repo.Save(ctx, event); err != nil {
		zap.L().Panic("Event", zap.Stringer("save", event.ID()), zap.Error(err))
	}

	for _, consumer := range b.consumers {
		consumer <- event
	}
}

// register регистрирует шаги бизнес-логики — источники команд, правила обработки доменных событий и потребителей
// событий.
func (b *Bus) register(step interface{}) {
	if consumer, ok := step.(domain.EventConsumer); ok {
		newChan := make(chan domain.Event)
		b.consumers = append(b.consumers, newChan)

		b.wg.Add(1)

		go func() {
			defer b.wg.Done()
			consumer.StartHandleEvent(b.loopCtx, newChan)
		}()
	}

	if source, ok := step.(domain.CommandSource); ok {
		b.wg.Add(1)

		go func() {
			defer b.wg.Done()
			source.StartProduceCommands(b.loopCtx, b.commands)
		}()
	}
}
