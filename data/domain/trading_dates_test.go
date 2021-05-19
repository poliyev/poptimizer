package domain

import (
	"context"
	"github.com/WLM1ke/gomoex"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTradingDatesFirstUpdate(t *testing.T) {
	table := TradingDates{iss: gomoex.NewISSClient(http.DefaultClient)}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(RowsReplaced)
	assert.True(t, ok)
}

func TestTradingDatesReplaceUpdate(t *testing.T) {
	rows := []gomoex.Date{{}}
	table := TradingDates{iss: gomoex.NewISSClient(http.DefaultClient), Rows: rows}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(RowsReplaced)
	assert.True(t, ok)
}

func TestTradingDatesEmptyUpdate(t *testing.T) {
	iss := gomoex.NewISSClient(http.DefaultClient)
	rows, _ := iss.MarketDates(context.Background(), gomoex.EngineStock, gomoex.MarketShares)
	table := TradingDates{iss: iss, Rows: rows}

	events := table.Update(context.Background())
	assert.Nil(t, events)
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

var testEvent = UpdateRequired{&TradingDates{ID: NewID(GroupTradingDates, GroupTradingDates), iss: nil}}

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
