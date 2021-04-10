package app

import (
	"context"
	"go.uber.org/zap"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"sync"
	"time"
)

type Bus struct {
	repo            *adapters.Repo
	handlersTimeout time.Duration

	commands  chan domain.Command
	events    chan domain.Event
	consumers []chan domain.Event

	wg          sync.WaitGroup
	loopCtx     context.Context
	loopCancel  context.CancelFunc
	loopStopped chan interface{}
}

func (b *Bus) Name() string {
	return "BUS"
}

func (b *Bus) Start(_ context.Context) error {
	go func() {
		b.loop(b.loopCtx)
		close(b.loopStopped)
	}()
	return nil
}

func (b *Bus) Shutdown(ctx context.Context) error {
	b.loopCancel()

	select {
	case <-b.loopStopped:
		return nil
	case <-ctx.Done():
		return context.DeadlineExceeded
	}
}

func (b *Bus) loop(ctx context.Context) {
	b.events = make(chan domain.Event)

	for {
		b.wg.Add(1)
		select {
		case cmd := <-b.commands:
			zap.L().Info(b.Name(), zap.Stringer("cmd", cmd.ID()))
			go func() {
				defer b.wg.Done()
				b.handleOneCommand(ctx, cmd)
			}()
		case event := <-b.events:
			zap.L().Info(b.Name(), zap.Stringer("event", event.ID()))
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
	ctx, cancel := context.WithTimeout(context.Background(), b.handlersTimeout)
	defer cancel()

	table, err := b.repo.Load(ctx, cmd.ID())
	if err != nil {
		zap.L().Panic(b.Name(), zap.Stringer("load", cmd.ID()), zap.Error(err))
	}
	for _, event := range table.HandleCommand(ctx, cmd) {
		b.events <- event
	}
}

func (b *Bus) handleOneEvent(ctx context.Context, event domain.Event) {
	ctx, cancel := context.WithTimeout(context.Background(), b.handlersTimeout)
	defer cancel()

	err := b.repo.Save(ctx, event)
	if err != nil {
		zap.L().Panic(b.Name(), zap.Stringer("save", event.ID()), zap.Error(err))
	}
	for _, consumer := range b.consumers {
		consumer <- event
	}
}

func (b *Bus) register(step interface{}) {
	if b.commands == nil {
		b.commands = make(chan domain.Command)
		ctx, loopCancel := context.WithCancel(context.Background())
		b.loopCtx = ctx
		b.loopCancel = loopCancel
		b.loopStopped = make(chan interface{})
	}

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
