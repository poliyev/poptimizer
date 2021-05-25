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

func TestDateToISSString(t *testing.T) {
	in := time.Date(2021, 5, 2, 0, 0, 0, 0, time.UTC)
	out := "2021-05-02"

	assert.Equal(t, out, dateToISSString(in))
}

type TestCandlesGateway struct {
	start string
	last  string
	rows  []gomoex.Candle
	err   error
}

func (t *TestCandlesGateway) MarketDates(_ context.Context, _, _ string) ([]gomoex.Date, error) {
	panic("implement me")
}

func (t *TestCandlesGateway) MarketCandles(_ context.Context, engine, market, ticker, start, last string,
	interval int) ([]gomoex.Candle, error) {
	if engine != "currency" || market != "selt" || ticker != "USD000UTSTOM" || start != t.start || last != t.
		last || interval != 24 {
		panic("некорректные аргументы")
	}

	return t.rows, t.err
}

func TestUSDGatewayError(t *testing.T) {
	iss := TestCandlesGateway{"", "2021-05-03", nil, fmt.Errorf("err")}
	last := time.Date(2021, 5, 3, 0, 0, 0, 0, time.UTC)
	dates := TradingDates{iss: &iss, Rows: []gomoex.Date{{Till: last}}}
	table := USD{iss: &iss, dates: &dates}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(UpdateError)
	assert.True(t, ok)
}

func TestUSDFirstUpdate(t *testing.T) {
	iss := TestCandlesGateway{"", "2021-06-04", []gomoex.Candle{{}}, nil}
	last := time.Date(2021, 6, 4, 0, 0, 0, 0, time.UTC)
	dates := TradingDates{iss: &iss, Rows: []gomoex.Date{{Till: last}}}
	table := USD{iss: &iss, dates: &dates}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(RowsAppended)
	assert.True(t, ok)
}

func TestUSDUpdateNotMatch(t *testing.T) {
	candle := gomoex.Candle{
		Begin:  time.Date(2019, 2, 5, 0, 0, 0, 0, time.UTC),
		End:    time.Now(),
		Open:   100.0,
		Close:  101.0,
		High:   102.0,
		Low:    98.0,
		Value:  11111.0,
		Volume: 70,
	}
	iss := TestCandlesGateway{"2019-02-05", "2021-07-14", []gomoex.Candle{{}, {}}, nil}
	last := time.Date(2021, 7, 14, 0, 0, 0, 0, time.UTC)
	dates := TradingDates{iss: &iss, Rows: []gomoex.Date{{Till: last}}}
	table := USD{iss: &iss, dates: &dates, Rows: []gomoex.Candle{{}, candle}}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	_, ok := events[0].(UpdateError)
	assert.True(t, ok)
}

func TestUSDUpdateMatchedUpdate(t *testing.T) {
	candle := gomoex.Candle{
		Begin:  time.Date(2018, 3, 5, 0, 0, 0, 0, time.UTC),
		End:    time.Now(),
		Open:   100.0,
		Close:  101.0,
		High:   102.0,
		Low:    98.0,
		Value:  11111.0,
		Volume: 70,
	}
	iss := TestCandlesGateway{"2018-03-05", "2021-07-14", []gomoex.Candle{candle, {}, {}, {}}, nil}
	last := time.Date(2021, 7, 14, 0, 0, 0, 0, time.UTC)
	dates := TradingDates{iss: &iss, Rows: []gomoex.Date{{Till: last}}}
	table := USD{iss: &iss, dates: &dates, Rows: []gomoex.Candle{{}, candle}}

	events := table.Update(context.Background())
	assert.Equal(t, 1, len(events))

	change, ok := events[0].(RowsAppended)
	assert.True(t, ok)
	assert.Equal(t, 3, len(change.Rows.([]gomoex.Candle)))
}

func TestUpdateUSDWrongEvent(t *testing.T) {
	in := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		assert.NotPanics(t, func() {
			producer := NewUpdateUSD(nil)
			producer.Activate(ctx, in, nil)
		})
	}()

	events := []Event{
		RowsReplaced{ID: NewID("a", GroupTradingDates)},
		RowsReplaced{ID: NewID(GroupTradingDates, "b")},
		UpdateError{ID: NewID(GroupTradingDates, GroupTradingDates)},
	}

	for _, event := range events {
		in <- event
	}

	cancel()
	wg.Wait()
}

func TestUpdateUSDRightEvent(t *testing.T) {
	in := make(chan Event)
	out := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		producer := UpdateUSD{}
		producer.Activate(ctx, in, out)
	}()

	go func() {
		in <- RowsReplaced{ID: NewID(GroupTradingDates, GroupTradingDates)}
	}()

	event := <-out
	_, ok := event.(UpdateRequired)
	assert.True(t, ok)

	cancel()
	wg.Wait()
}

func TestMakeUSDEvent(t *testing.T) {
	event := makeUSDEvent(nil)

	assert.Equal(t, 2, len(event.Templates))

	first := event.Templates[0].(*USD)
	second := event.Templates[1].(*TradingDates)

	assert.Equal(t, Group(GroupUSD), first.Group())
	assert.Equal(t, Group(GroupTradingDates), second.Group())

	assert.Equal(t, Name(GroupUSD), first.Name())
	assert.Equal(t, Name(GroupTradingDates), second.Name())

	assert.Equal(t, second, first.dates)
}
