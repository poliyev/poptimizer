package app

import (
	"context"
	"poptimizer/data/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestEventConsumer struct {
	flow chan<- string
}

func (t TestEventConsumer) StartHandleEvent(ctx context.Context, source <-chan domain.Event) {
	for {
		select {
		case event := <-source:
			t.flow <- "Event consumed " + event.ID().String()
		case <-ctx.Done():
			t.flow <- "Stopped"

			return
		}
	}
}

type TestCommandSource struct {
	flow chan<- string
}

func (t TestCommandSource) StartProduceCommands(ctx context.Context, output chan<- domain.Command) {
	cmd := domain.NewTableID("cg", "cn")
	output <- cmd
	t.flow <- "Command send"

	<-ctx.Done()

	t.flow <- "Stopped"
}

type TestTable struct {
	domain.TableID
	flow chan<- string
}

func (t TestTable) ID() domain.TableID {
	return t.TableID
}

func (t TestTable) HandleCommand(_ context.Context, cmd domain.Command) []domain.Event {
	t.flow <- "TableHandles " + cmd.ID().String()

	event := domain.NewTableID("eg", "en")

	return []domain.Event{event}
}

type TestRepo struct {
	flow chan<- string
}

func (t TestRepo) Load(_ context.Context, id domain.TableID) (domain.Table, error) {
	t.flow <- "TableLoad " + id.String()

	return TestTable{TableID: id, flow: t.flow}, nil
}

func (t TestRepo) Save(_ context.Context, event domain.Event) error {
	t.flow <- "Event saved " + event.ID().String()

	return nil
}

func TestBusFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	flow := make(chan string, 3)

	bus := Bus{
		repo:             TestRepo{flow},
		handlersTimeouts: time.Second,

		commands: make(chan domain.Command),
		events:   make(chan domain.Event),

		loopCtx:     ctx,
		loopCancel:  cancel,
		loopStopped: make(chan struct{}),
	}

	bus.register(TestEventConsumer{flow})
	bus.register(TestEventConsumer{flow})
	bus.register(TestCommandSource{flow})

	if bus.Start(context.Background()) != nil {
		t.Error("Не удалось запустить шину")
	}

	assert.Equal(t, "Command send", <-flow)
	assert.Equal(t, "TableLoad ID(cg, cn)", <-flow)
	assert.Equal(t, "TableHandles ID(cg, cn)", <-flow)
	assert.Equal(t, "Event saved ID(eg, en)", <-flow)
	assert.Equal(t, "Event consumed ID(eg, en)", <-flow)
	assert.Equal(t, "Event consumed ID(eg, en)", <-flow)

	if bus.Shutdown(context.Background()) != nil {
		t.Error("Не удалось остановить шину")
	}

	assert.Equal(t, "Stopped", <-flow)
	assert.Equal(t, "Stopped", <-flow)
	assert.Equal(t, "Stopped", <-flow)
}

func TestBusName(t *testing.T) {
	bus := NewBus(TestRepo{}, time.Second)

	assert.Equal(t, "Bus", bus.Name())
}
