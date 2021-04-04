package domain

import (
	"github.com/WLM1ke/gomoex"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewMainFactory(t *testing.T) {
	factory := NewMainFactory(gomoex.NewISSClient(http.DefaultClient))
	table := factory.NewTable(groupTradingDates, groupTradingDates)
	_, ok := table.(*TradingDates)
	assert.True(t, ok)
}

func TestNewMainPanicOnWrongGroup(t *testing.T) {
	factory := NewMainFactory(gomoex.NewISSClient(http.DefaultClient))
	assert.Panics(t, func() {
		factory.NewTable("Bad", "Bad")
	})
}

func TestNewMainPanicOnWrongSingleton(t *testing.T) {
	factory := NewMainFactory(gomoex.NewISSClient(http.DefaultClient))
	assert.Panics(t, func() {
		factory.NewTable(groupTradingDates, "Bad")
	})
}
