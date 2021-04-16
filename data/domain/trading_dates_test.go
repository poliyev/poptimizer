package domain

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/WLM1ke/gomoex"
	"github.com/stretchr/testify/assert"
)

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

func TestTradingDayAppStart(t *testing.T) {
	out := UpdateTable{TableID{GroupTradingDates, GroupTradingDates}}

	output := make(chan Command)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		producer := NewCheckTradingDay()
		producer.StartProduceCommands(ctx, output)
	}()
	assert.Equal(t, &out, <-output)

	cancel()
	wg.Wait()
}

func TestTradingDayNextUpdate(t *testing.T) {
	out := UpdateTable{TableID{GroupTradingDates, GroupTradingDates}}

	timer := make(chan time.Time)

	output := make(chan Command)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		producer := CheckTradingDay{ticker: timer, stopFn: func() {}, tz: testMoscowTZ}
		producer.StartProduceCommands(ctx, output)
	}()

	<-output

	// После публикации данных на ISS должна отправляться команда
	now := time.Now()
	timer <- nextISSDailyUpdate(now, testMoscowTZ).Add(time.Second)

	assert.Equal(t, &out, <-output)

	// До начала следующего дня обновление таймера не порождает команд
	close(output)
	timer <- nextISSDailyUpdate(now, testMoscowTZ).Add(time.Hour * 24)

	cancel()
	wg.Wait()
}

func TestTradingDatesFirstUpdate(t *testing.T) {
	table := TradingDates{iss: gomoex.NewISSClient(http.DefaultClient)}
	cmd := UpdateTable{TableID{GroupTradingDates, GroupTradingDates}}

	events := table.HandleCommand(context.Background(), &cmd)
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(RowsReplaced)
	assert.True(t, ok)
}

func TestTradingDatesReplaceUpdate(t *testing.T) {
	rows := []gomoex.Date{{}}
	table := TradingDates{iss: gomoex.NewISSClient(http.DefaultClient), Rows: rows}
	cmd := UpdateTable{TableID{GroupTradingDates, GroupTradingDates}}

	events := table.HandleCommand(context.Background(), &cmd)
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(RowsReplaced)
	assert.True(t, ok)
}

func TestTradingDatesEmptyUpdate(t *testing.T) {
	iss := gomoex.NewISSClient(http.DefaultClient)
	rows, _ := iss.MarketDates(context.Background(), gomoex.EngineStock, gomoex.MarketShares)
	table := TradingDates{iss: iss, Rows: rows}
	cmd := UpdateTable{TableID{GroupTradingDates, GroupTradingDates}}

	events := table.HandleCommand(context.Background(), &cmd)
	assert.Nil(t, events)
}
