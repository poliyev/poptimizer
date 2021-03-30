package tables

import (
	"time"
)

type Event interface {
	Day() time.Time
}

type EventHandler interface {
	Match(event Event) bool
	Handle(event Event)
}

type (
	Group   string
	Name    string
	Command interface {
	}
)

type UpdateTable struct {
	Group Group
	Name  Name
	Day   time.Time
}

// InitTableFactory создает фабрику и регистрирует все шаблоны таблиц.
//func InitTableFactory() Factory {
//	factory := MainFactory{}
//	iss := gomoex.NewISSClient(http.DefaultClient)
//	factory.registerTemplate("trading_dates", TradingDatesFactory{iss}, true)
//	return &factory
//}
