package app

import (
	"context"
	"fmt"
	"poptimizer/data/domain"
	"testing"

	"github.com/stretchr/testify/assert"
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
			t.flow <- fmt.Sprintf("Event consumed ID(%s, %s)", event.Group(), event.Name())
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
