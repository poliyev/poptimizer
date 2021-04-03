package domain

import (
	"github.com/WLM1ke/gomoex"
)

// MainFactory регистрирует другие фабрики и создает любые таблицы.
type MainFactory struct {
	groupTemplates map[Group]Factory
	singletons     map[Group]bool
}

// Каждая фабрика создает талицы из одной группы. В группе может быть единственная талица,
// тогда имя таблицы должно совпадать с именем группы.
func (t *MainFactory) registerGroupFactory(group Group, factory Factory, singleton bool) {
	if t.groupTemplates == nil {
		t.groupTemplates = make(map[Group]Factory, 0)
		t.singletons = make(map[Group]bool, 0)
	}
	t.groupTemplates[group] = factory
	t.singletons[group] = singleton
}

// NewTable - создает таблицу и проверяет, что указано корректное имя таблицы для групп с одной таблицей.
func (t *MainFactory) NewTable(group Group, name Name) Table {
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
	factory := MainFactory{map[Group]Factory{}, map[Group]bool{}}

	factory.registerGroupFactory(groupTradingDates, tradingDatesFactory{iss}, true)

	return &factory
}
