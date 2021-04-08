package domain

import (
	"github.com/WLM1ke/gomoex"
)

// mainFactory создает любые таблицы.
type mainFactory struct {
	iss *gomoex.ISSClient
}

// NewTable - создает таблицу и проверяет, что указано корректное имя таблицы для групп с одной таблицей.
func (t *mainFactory) NewTable(id TableID) Table {
	switch {
	case id.Group == GroupTradingDates && id.Name == GroupTradingDates:
		return &TradingDates{TableID: id, iss: t.iss}
	default:
		panic("Некорректное ID таблицы")
	}
}

// NewMainFactory - создает главную фабрику и регистрирует все доступные группы таблиц.
func NewMainFactory(iss *gomoex.ISSClient) Factory {
	factory := mainFactory{iss}

	return &factory
}
