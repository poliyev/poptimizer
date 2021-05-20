package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"poptimizer/data/domain"
	"testing"
)

type TestRule struct {
	flow    chan<- string
	atStart bool
}

func (t TestRule) Activate(ctx context.Context, in <-chan domain.Event, out chan<- domain.Event) {
	if t.atStart {
		out <- domain.NewID("g", "n")
	}

	for {
		select {
		case event := <-in:
			t.flow <- "Event consumed " + event.String()
		case <-ctx.Done():
			t.flow <- "Stopped"

			return
		}
	}
}

func TestBusFlow(t *testing.T) {
	flow := make(chan string, 2)
	rules := []domain.Rule{TestRule{flow, false}, TestRule{flow, true}}

	bus := Bus{rules: rules}

	if bus.Start(context.Background()) != nil {
		t.Error("Не удалось запустить шину")
	}

	assert.Equal(t, "Event consumed ID(g, n)", <-flow)
	assert.Equal(t, "Event consumed ID(g, n)", <-flow)

	if bus.Shutdown(context.Background()) != nil {
		t.Error("Не удалось остановить шину")
	}

	assert.Equal(t, "Stopped", <-flow)
	assert.Equal(t, "Stopped", <-flow)
}
