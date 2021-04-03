package domain

import (
	"context"
	"errors"
	"github.com/WLM1ke/gomoex"
)

var ErrRowsValidationErr = errors.New("ошибка валидации данных")

// groupTradingDates - группа таблицы с торговыми данными.
const groupTradingDates = "trading_dates"

// TradingDates - таблица с диапазоном торговых дат.
type TradingDates struct {
	ID

	iss *gomoex.ISSClient

	Rows []gomoex.Date
}

// HandleCommand - полностью переписывает таблицу, если появились новые данные о торговых датах.
func (t *TradingDates) HandleCommand(ctx context.Context, _ Command) []Event {
	newRows, err := t.iss.MarketDates(ctx, gomoex.EngineStock, gomoex.MarketShares)

	switch {
	case err != nil:
		panic("не удалось получить данные от ISS")
	case len(newRows) != 1:
		panic("ошибка валидации данных от ISS")
	case t.Rows == nil, !newRows[0].Till.Equal(t.Rows[0].Till):
		return []Event{&RowsReplaced{t.ID, newRows}}
	default:
		return nil
	}
}

type tradingDatesFactory struct {
	iss *gomoex.ISSClient
}

func (t tradingDatesFactory) NewTable(group Group, name Name) Table {

	return &TradingDates{ID: ID{group, name}, iss: t.iss}
}
