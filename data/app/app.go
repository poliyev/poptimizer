package app

import (
	"context"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
)

const (
	MongoURI = "mongodb://localhost:27017"
	mongoDB  = "new_data"
)

type App struct {
	repo *adapters.Repo
}

func (a *App) Run(ctx context.Context) {
	if a.repo == nil {
		iss := adapters.NewISSClient()
		factory := domain.NewMainFactory(iss)
		// TODO: закрывать РЕПО
		a.repo = adapters.NewRepo(ctx, MongoURI, mongoDB, factory)
	}

	bus := Bus{repo: a.repo}

	steps := []interface{}{
		// Источники команд
		&domain.CheckTradingDay{},
		// Правила

		// Потребители сообщений
	}
	for _, step := range steps {
		bus.register(step)
	}

	bus.Run(ctx)
}

func (a App) GetJson(ctx context.Context) ([]byte, error) {
	return a.repo.ViewJOSN(ctx, domain.TableID{"trading_dates", "trading_dates"})
}
