package tables

import (
	"context"
	"github.com/WLM1ke/gomoex"
	"time"
)

const tradingDates = "trading_dates"

var moexZone = loadMoexZone()

func loadMoexZone() *time.Location {
	zone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("не удалось инициализировать информацию о часовом поясе")
	}
	return zone
}

type CheckTradingDates struct {
}

func (c CheckTradingDates) Group() TableGroup {
	return tradingDates
}

func (c CheckTradingDates) Name() TableName {
	return tradingDates
}

func (c CheckTradingDates) Type() CommandType {
	return tradingDates
}

type TradingDates struct {
	timestamp time.Time
	events    []Event
	iss       *gomoex.ISSClient
	rows      []Row
}

func NewTradingDates(iss *gomoex.ISSClient) TradingDates {
	return TradingDates{iss: iss}
}

func (t *TradingDates) Group() TableGroup {
	return tradingDates
}

func (t *TradingDates) Name() TableName {
	return tradingDates
}

func (t *TradingDates) Timestamp() time.Time {
	return t.timestamp
}

func (t *TradingDates) FlushEvents() []Event {
	events := t.events
	t.events = make([]Event, 0)
	return events
}

func (t *TradingDates) Handle(ctx context.Context, command Command) error {
	if command.Type() != tradingDates {
		return ErrWrongCommandType
	}

	if t.updateCond() {
		rows, err := t.prepareRows(ctx, command)
		if err != nil {
			return err
		}

		err = t.validateRows(rows)
		if err != nil {
			return err
		}

		t.addEvent(rows)
	}
	return nil
}

func (t *TradingDates) updateCond() bool {
	now := time.Now().In(moexZone)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 45, 0, 0, moexZone)
	if end.After(now) {
		end = end.AddDate(0, 0, -1)
	}
	return end.After(t.timestamp)
}

func (t *TradingDates) prepareRows(ctx context.Context, command Command) ([]Row, error) {
	rawRows, err := t.iss.MarketDates(ctx, gomoex.EngineStock, gomoex.MarketShares)
	if err != nil {
		return nil, err
	}

	rows := make([]Row, len(rawRows))
	for n, row := range rawRows {
		rows[n] = Row(row)
	}
	return rows, nil
}

func (t *TradingDates) validateRows(rows []Row) error {
	if len(rows) != 1 {
		return ErrRowsValidationErr
	}

	return nil
}

func (t *TradingDates) addEvent(rows []Row) {
	t.timestamp = time.Now()
	if t.replaceRows() {
		t.rows = rows
	} else {
		t.rows = append(t.rows, rows...)
	}

	newEvent := Event{
		Group:       tradingDates,
		Name:        tradingDates,
		Timestamp:   t.timestamp,
		ReplaceRows: t.replaceRows(),
		NewRows:     rows,
	}
	t.events = append(t.events, newEvent)
}

func (t *TradingDates) replaceRows() bool {
	return true
}
