package app

import (
	"context"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
)

type App struct {
	repo *adapters.Repo
}

func (a *App) Run(ctx context.Context) {
	if a.repo == nil {
		iss := adapters.NewISSClient()
		factory := domain.NewMainFactory(iss)
		a.repo = adapters.NewRepo(factory)
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
	return a.repo.ViewJOSN(ctx, "trading_dates", "trading_dates")
}
