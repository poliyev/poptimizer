package tables

import (
	"github.com/WLM1ke/gomoex"
	"net/http"
)

// InitTableFactory создает фабрику и регистрирует все шаблоны таблиц.
func InitTableFactory() *TableFactory {
	factory := TableFactory{}
	iss := gomoex.NewISSClient(http.DefaultClient)
	factory.registerTemplate(TradingDatesFactory{iss})
	return &factory
}
