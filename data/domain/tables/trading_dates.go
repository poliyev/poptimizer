package tables

import (
	"context"
	"errors"
	"github.com/WLM1ke/gomoex"
	"poptimizer/data/domain"
)

var ErrRowsValidationErr = errors.New("ошибка валидации данных")

// TradingDates - таблица с диапазоном торговых дат.
//
// ID таблицы должна заполнять фабрика.
// Ряды таблицы и последняя торговая дата должны грузиться из базы.
type TradingDates struct {
	group domain.Group
	name  domain.Name

	iss *gomoex.ISSClient

	Rows []gomoex.Date
}

func (t *TradingDates) Group() domain.Group {
	return t.group
}

func (t *TradingDates) Name() domain.Name {
	return t.name
}

func (t *TradingDates) Update(ctx context.Context, _ domain.Command) (domain.Event, error) {
	newRows, err := t.iss.MarketDates(ctx, gomoex.EngineStock, gomoex.MarketShares)

	switch {
	case err != nil:
		return domain.Event{}, err
	case len(newRows) != 1:
		return domain.Event{}, ErrRowsValidationErr
	case t.Rows == nil, !newRows[0].Till.Equal(t.Rows[0].Till):
		event := domain.Event{t.group, t.name, true, newRows}
		return event, nil
	default:
		return domain.Event{}, nil
	}
}

type TradingDatesFactory struct {
	iss *gomoex.ISSClient
}

func (t TradingDatesFactory) NewTable(group domain.Group, name domain.Name) domain.Table {
	return &TradingDates{group: group, name: name, iss: t.iss}
}
