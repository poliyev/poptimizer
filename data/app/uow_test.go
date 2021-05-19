package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"poptimizer/data/domain"
	"sync"
	"testing"
	"time"
)

type TestTable struct {
	domain.ID
	outEvent domain.Event
}

func (t *TestTable) Update(_ context.Context) []domain.Event {
	return []domain.Event{t.outEvent}
}

type TestRepo struct {
}

func (t TestRepo) Unmarshal(_ context.Context, event domain.UpdateRequired) (domain.Table, error) {
	return event.Template, nil
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

		uow.Activate(ctx, in, nil)
	}()

	// Была ошибка с завершением работы после первого необрабатываемого события
	in <- event
	in <- event

	wg.Wait()
	cancel()
}

func TestUoWAppendRowsUpdate(t *testing.T) {
	table := TestTable{domain.ID{}, domain.RowsAppended{}}
	event := domain.UpdateRequired{Template: &table}

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
	event := domain.UpdateRequired{Template: &table}

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
