package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"poptimizer/data/domain"
	"sync"
	"testing"
)

func TestNonErrorEvent(t *testing.T) {
	rule := ErrorsHandler{}
	in := make(chan domain.Event)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		assert.NotPanics(t, func() {
			rule.Activate(ctx, in, nil)
		},
			"Правило должно игнорировать не ошибочные события.",
		)
	}()

	in <- domain.RowsAppended{}

	cancel()
	wg.Wait()
}

func TestErrorEvent(t *testing.T) {
	rule := ErrorsHandler{}
	in := make(chan domain.Event)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		assert.Panics(t, func() {
			rule.Activate(ctx, in, nil)
		},
			"Правило должно паниковать при ошибочных событиях",
		)
	}()

	in <- domain.UpdateError{}

	cancel()
	wg.Wait()
}
