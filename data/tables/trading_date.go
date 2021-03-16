package tables

import (
	"context"
	"github.com/WLM1ke/gomoex"
	"time"
)

var moexZone = loadMoexZone()

func loadMoexZone() *time.Location {
	zone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("не удалось инициализировать информацию о часовом поясе")
	}
	return zone
}

type TradingDates struct {
	iss     *gomoex.ISSClient
	rows    []gomoex.Date
	newRows []gomoex.Date
}

func (t *TradingDates) replaceRows() bool {
	return true
}

func (t *TradingDates) updateCond(timestamp time.Time) bool {
	now := time.Now().In(moexZone)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 45, 0, 0, moexZone)
	if end.After(now) {
		end = end.AddDate(0, 0, -1)
	}
	return end.After(timestamp)
}

func (t *TradingDates) prepareRows(ctx context.Context, _ Command) (err error) {
	t.newRows, err = t.iss.MarketDates(ctx, gomoex.EngineStock, gomoex.MarketShares)
	if err != nil {
		return err
	}

	return nil
}

func (t *TradingDates) validateRows() error {
	if len(t.newRows) != 1 {
		return ErrRowsValidationErr
	}

	return nil
}

func (t *TradingDates) addNewRows() []Row {
	t.rows = t.newRows

	genRows := make([]Row, len(t.newRows))
	for n, row := range t.newRows {
		genRows[n] = Row(row)
	}

	return genRows
}

func (t *TradingDates) replace() bool {
	return true
}

type TradingDatesFactory struct {
	iss *gomoex.ISSClient
}

func (t TradingDatesFactory) group() Group {
	return "trading_dates"
}

func (t TradingDatesFactory) singleton() bool {
	return true
}

func (t TradingDatesFactory) NewTable(name Name) Table {
	return &BaseTable{ID: ID{"trading_dates", name}, tableTemplate: &TradingDates{iss: t.iss}}
}
