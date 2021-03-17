package tables

import (
	"github.com/WLM1ke/gomoex"
	"net/http"
)

// InitTableFactory создает фабрику и регистрирует все шаблоны таблиц.
func InitTableFactory() Factory {
	factory := MainFactory{}
	iss := gomoex.NewISSClient(http.DefaultClient)
	factory.registerTemplate("trading_dates", TradingDatesFactory{iss}, true)
	return &factory
}
