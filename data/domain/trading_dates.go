package domain

import (
	"context"
	"errors"
	"github.com/WLM1ke/gomoex"
)

var ErrRowsValidationErr = errors.New("ошибка валидации данных")

// TradingDates - таблица с диапазоном торговых дат.
//
// ID таблицы должна заполнять фабрика.
// Ряды таблицы и последняя торговая дата должны грузиться из базы.
type TradingDates struct {
	ID
	events []Event

	iss *gomoex.ISSClient

	Rows []gomoex.Date
}

func (t *TradingDates) Events() []Event {
	events := t.events
	t.events = nil
	return events
}

func (t *TradingDates) Update(ctx context.Context, _ Command) {
	newRows, err := t.iss.MarketDates(ctx, gomoex.EngineStock, gomoex.MarketShares)

	var event Event
	switch {
	case err != nil:
		event = &TableUpdateErrOccurred{t.ID, err}
	case len(newRows) != 1:
		event = &TableUpdateErrOccurred{t.ID, ErrRowsValidationErr}
	case t.Rows == nil, !newRows[0].Till.Equal(t.Rows[0].Till):
		event = &RowsReplaced{t.ID, newRows}
	default:
		return
	}

	t.events = append(t.events, event)
}

type TradingDatesFactory struct {
	iss *gomoex.ISSClient
}

func (t TradingDatesFactory) NewTable(group Group, name Name) Table {

	return &TradingDates{ID: ID{group, name}, iss: t.iss}
}
