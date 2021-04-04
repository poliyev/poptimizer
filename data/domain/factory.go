package domain

import (
	"github.com/WLM1ke/gomoex"
)

// mainFactory регистрирует другие фабрики и создает любые таблицы.
type mainFactory struct {
	groupTemplates map[Group]Factory
	singletons     map[Group]bool
}

// Каждая фабрика создает талицы из одной группы. В группе может быть единственная талица,
// тогда имя таблицы должно совпадать с именем группы.
func (t *mainFactory) registerGroupFactory(group Group, factory Factory, singleton bool) {
	t.groupTemplates[group] = factory
	t.singletons[group] = singleton
}

// NewTable - создает таблицу и проверяет, что указано корректное имя таблицы для групп с одной таблицей.
func (t *mainFactory) NewTable(group Group, name Name) Table {
	factory, ok := t.groupTemplates[group]
	if !ok {
		panic("Незарегистрированная группа")
	}

	if t.singletons[group] && group != Group(name) {
		panic("Некорректное имя таблицы")
	}

	return factory.NewTable(group, name)
}

// NewMainFactory - создает главную фабрику и регистрирует все доступные группы таблиц.
func NewMainFactory(iss *gomoex.ISSClient) Factory {
	factory := mainFactory{map[Group]Factory{}, map[Group]bool{}}

	factory.registerGroupFactory(groupTradingDates, tradingDatesFactory{iss}, true)

	return &factory
}
