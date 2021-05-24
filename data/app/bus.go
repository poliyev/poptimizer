package app

import (
	"context"
	"fmt"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"sync"
	"time"

	"github.com/WLM1ke/gomoex"
	"go.uber.org/zap"
)

// Bus - шина для обработки событий, поддерживающую интерфейс модуля приложения.
//
// Направляет входящий поток событий набору бизнес-правил, читает их исходящие события и закольцовывает их обратно на вход правил.
type Bus struct {
	rules []domain.Rule

	wg         sync.WaitGroup
	loopCancel context.CancelFunc
}

// NewBus - создает шину для обработки событий.
func NewBus(repo Repo, eventBusTimeouts time.Duration, iss *gomoex.ISSClient) *Bus {
	rules := []domain.Rule{
		ErrorsHandler{},
		NewUoW(repo, eventBusTimeouts),
		domain.NewUpdateTradingDates(iss),
		domain.NewUpdateUSD(iss),
	}

	bus := Bus{
		rules: rules,
	}

	return &bus
}

// Start - запускает работу правил и основного цикла обработки событий.
func (b *Bus) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	b.loopCancel = cancel

	events := make(chan domain.Event)
	consumers := make([]chan domain.Event, 0)

	for _, rule := range b.rules {
		rule := rule
		in := make(chan domain.Event)
		consumers = append(consumers, in)

		b.wg.Add(1)

		go func() {
			defer b.wg.Done()

			zap.L().Info("Activated", adapters.TypeField("rule", rule))
			rule.Activate(ctx, in, events)
			zap.L().Info("Deactivated", adapters.TypeField("rule", rule))
		}()
	}

	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		b.loop(ctx, events, consumers)
	}()

	return nil
}

func (b *Bus) loop(ctx context.Context, events chan domain.Event, consumers []chan domain.Event) {
	for {
		select {
		case event := <-events:
			b.wg.Add(1)

			go func() {
				defer b.wg.Done()

				zap.L().Info("Handling", adapters.EventField(event))

				for _, consumer := range consumers {
					consumer <- event
				}
			}()

		case <-ctx.Done():
			return
		}
	}
}

// Shutdown - завершает работу основного цикла обработки команд и событий.
func (b *Bus) Shutdown(ctx context.Context) error {
	b.loopCancel()

	select {
	case <-b.waitCancel():
		return nil
	case <-ctx.Done():
		return fmt.Errorf("bus shutdown error: %w", ctx.Err())
	}
}

func (b *Bus) waitCancel() chan struct{} {
	stopped := make(chan struct{})

	go func() {
		b.wg.Wait()
		close(stopped)
	}()

	return stopped
}
