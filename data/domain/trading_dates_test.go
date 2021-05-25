package domain

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/WLM1ke/gomoex"

	"github.com/stretchr/testify/assert"
)

type TestTradingDatesGateway struct {
	rows []gomoex.Date
	err  error
}

func (t *TestTradingDatesGateway) MarketDates(_ context.Context, engine string, market string) ([]gomoex.Date, error) {
	if engine != gomoex.EngineStock || market != gomoex.MarketShares {
		panic("некорректные аргументы")
	}

	return t.rows, t.err
}

func TestTradingDatesGatewayError(t *testing.T) {
	iss := TestTradingDatesGateway{nil, fmt.Errorf("err")}
	table := TradingDates{iss: &iss}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(UpdateError)
	assert.True(t, ok)
}

func TestTradingDatesValidationError(t *testing.T) {
	rows := []gomoex.Date{{}, {}}
	iss := TestTradingDatesGateway{rows, nil}
	table := TradingDates{iss: &iss}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(UpdateError)
	assert.True(t, ok)
}

func TestTradingDatesFirstUpdate(t *testing.T) {
	rows := []gomoex.Date{{}}
	iss := TestTradingDatesGateway{rows, nil}
	table := TradingDates{iss: &iss}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(RowsReplaced)
	assert.True(t, ok)
}

func TestTradingDatesReplaceUpdate(t *testing.T) {
	rows := []gomoex.Date{{Till: time.Now()}}
	iss := TestTradingDatesGateway{rows, nil}
	table := TradingDates{iss: &iss, Rows: []gomoex.Date{{}}}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(RowsReplaced)
	assert.True(t, ok)
}

func TestTradingDatesNoUpdate(t *testing.T) {
	rows := []gomoex.Date{{}}
	iss := TestTradingDatesGateway{rows, nil}
	table := TradingDates{iss: &iss, Rows: []gomoex.Date{{}}}

	events := table.Update(context.Background())
	assert.Nil(t, events)
}

func TestTradingDatesLastDate(t *testing.T) {
	last := time.Now()
	table := TradingDates{Rows: []gomoex.Date{{Till: last}}}

	assert.Equal(t, last, table.lastDate())
}

var testMoscowTZ = prepareZone("Europe/Moscow")

func TestZone(t *testing.T) {
	moscow, _ := time.LoadLocation("Europe/Moscow")
	assert.Equal(t, testMoscowTZ, moscow)
}

func TestZonePanic(t *testing.T) {
	assert.Panics(t, func() { prepareZone("WrongZone") })
}

func TestBeforeNextISSDailyUpdate(t *testing.T) {
	in := time.Date(2021, 4, 3, 21, 44, 0, 0, time.UTC)
	out := time.Date(2021, 4, 4, 0, 45, 0, 0, testMoscowTZ)

	assert.Equal(t, out, nextISSDailyUpdate(in, testMoscowTZ))
}

func TestAfterNextISSDailyUpdate(t *testing.T) {
	in := time.Date(2021, 4, 3, 21, 46, 0, 0, time.UTC)
	out := time.Date(2021, 4, 5, 0, 45, 0, 0, testMoscowTZ)

	assert.Equal(t, out, nextISSDailyUpdate(in, prepareZone("Europe/Moscow")))
}

var testEvent = UpdateRequired{[]Table{&TradingDates{ID: NewID(GroupTradingDates, GroupTradingDates), iss: nil}}}

func TestTradingDayAppStartOutInNotBlocks(t *testing.T) {
	in := make(chan Event)
	out := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		rule := NewUpdateTradingDates(nil)
		rule.Activate(ctx, in, out)
	}()

	var firstOut, inNotBlocks bool

	for {
		select {
		case result := <-out:
			assert.Equal(t, testEvent, result)

			firstOut = true
		case in <- testEvent:
			inNotBlocks = true
		}

		if firstOut && inNotBlocks {
			cancel()
			wg.Wait()

			return
		}
	}
}

func TestTradingDayNextUpdate(t *testing.T) {
	timer := make(chan time.Time)

	out := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		producer := UpdateTradingDates{ticker: timer, stopFn: func() {}, tz: testMoscowTZ}
		producer.Activate(ctx, nil, out)
	}()

	<-out

	// После публикации данных на ISS должна отправляться команда
	now := time.Now()
	timer <- nextISSDailyUpdate(now, testMoscowTZ).Add(time.Second)

	assert.Equal(t, testEvent, <-out)

	// До начала следующего дня обновление таймера не порождает команд
	close(out)
	timer <- nextISSDailyUpdate(now, testMoscowTZ).Add(time.Hour * 24)

	cancel()
	wg.Wait()
}
