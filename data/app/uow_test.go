package app

import (
	"context"
	"poptimizer/data/domain"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestTable struct {
	domain.ID
	outEvent domain.Event
}

func (t *TestTable) Update(_ context.Context) []domain.Event {
	return []domain.Event{t.outEvent}
}

type TestRepo struct{}

func (t TestRepo) Unmarshal(_ context.Context, _ domain.Table) error {
	return nil
}

func (t TestRepo) Replace(_ context.Context, _ domain.RowsReplaced) error {
	return nil
}

func (t TestRepo) Append(_ context.Context, _ domain.RowsAppended) error {
	return nil
}

func TestUoWNotHandelNonUpdateEvents(t *testing.T) {
	event := domain.RowsAppended{}
	repo := TestRepo{}
	uow := NewUoW(&repo, time.Second)

	in := make(chan domain.Event)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		assert.NotPanics(t, func() {
			uow.Activate(ctx, in, nil)
		},
			"Правило должно паниковать при ошибочных событиях",
		)
	}()

	// Была ошибка с завершением работы после первого необрабатываемого события
	in <- event
	in <- event

	wg.Wait()
	cancel()
}

func TestUoWAppendRowsUpdate(t *testing.T) {
	table := TestTable{domain.ID{}, domain.RowsAppended{}}
	event := domain.UpdateRequired{Templates: []domain.Table{&table}}

	repo := TestRepo{}
	uow := NewUoW(&repo, time.Second)

	in := make(chan domain.Event)
	out := make(chan domain.Event)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		uow.Activate(ctx, in, out)
	}()

	go func() {
		in <- event
	}()

	_, ok := (<-out).(domain.RowsAppended)

	assert.True(t, ok)

	cancel()
	wg.Wait()
}

func TestUoWReplaceRowsUpdate(t *testing.T) {
	table := TestTable{domain.ID{}, domain.RowsReplaced{}}
	event := domain.UpdateRequired{Templates: []domain.Table{&table}}

	repo := TestRepo{}
	uow := NewUoW(&repo, time.Second)

	in := make(chan domain.Event)
	out := make(chan domain.Event)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		uow.Activate(ctx, in, out)
	}()

	go func() {
		in <- event
	}()

	_, ok := (<-out).(domain.RowsReplaced)

	assert.True(t, ok)

	cancel()
	wg.Wait()
}
