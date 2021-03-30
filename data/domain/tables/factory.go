package tables

import (
	"github.com/WLM1ke/gomoex"
	"poptimizer/data/domain"
)

// MainFactory регистрирует и создает таблицы.
type MainFactory struct {
	groupTemplates map[domain.Group]domain.Factory
	singletons     map[domain.Group]bool
}

func (t *MainFactory) registerTemplate(group domain.Group, factory domain.Factory, singleton bool) {
	if t.groupTemplates == nil {
		t.groupTemplates = make(map[domain.Group]domain.Factory, 0)
		t.singletons = make(map[domain.Group]bool, 0)
	}
	t.groupTemplates[group] = factory
	t.singletons[group] = singleton
}

func (t *MainFactory) NewTable(group domain.Group, name domain.Name) domain.Table {
	factory, ok := t.groupTemplates[group]
	if !ok {
		panic("незарегестрированая группа")
	}

	if t.singletons[group] && group != domain.Group(name) {
		panic("некорректное имя таблицы")
	}

	return factory.NewTable(group, name)
}

const (
	GroupTradingDates = "trading_dates"
)

func NewMainFactory(iss *gomoex.ISSClient) domain.Factory {
	factory := MainFactory{map[domain.Group]domain.Factory{}, map[domain.Group]bool{}}

	factory.registerTemplate(GroupTradingDates, TradingDatesFactory{iss}, true)

	return &factory
}
