package cfg

import (
	"context"
	"go.uber.org/zap"
	"poptimizer/data/adapters"
	"poptimizer/data/app"
	"poptimizer/data/domain"
	"poptimizer/data/ports"
	"time"
)

type Config struct {
	StartTimeout     time.Duration
	ShutdownTimeout  time.Duration
	RequestTimeout   time.Duration
	EventBusTimeouts time.Duration
	ServerAddr       string
	ISSMaxCons       int
	MongoURI         string
	MongoDB          string
}

type App struct {
	startTimeout    time.Duration
	shutdownTimeout time.Duration
	modules         []Module
}

func NewApp(cfg Config) *App {
	iss := adapters.NewISSClient(cfg.ISSMaxCons)
	factory := domain.NewMainFactory(iss)
	repo := adapters.NewRepo(cfg.MongoURI, cfg.MongoDB, factory)
	bus := app.NewBus(repo, cfg.EventBusTimeouts)

	modules := []Module{
		adapters.NewLogger(),
		repo,
		bus,
		ports.NewServer(cfg.ServerAddr, cfg.RequestTimeout, repo),
	}

	return &App{
		startTimeout:    cfg.StartTimeout,
		shutdownTimeout: cfg.ShutdownTimeout,
		modules:         modules,
	}
}

func (a *App) Run() {
	defer func() {
		zap.L().Info("App", zap.String("status", "stopped"))

	}()

	appCtx := appTerminationCtx()

	startCtx, startCancel := context.WithTimeout(appCtx, a.startTimeout)
	defer startCancel()

	for _, module := range a.modules {
		if err := module.Start(startCtx); err != nil {
			zap.L().Panic(module.Name(), zap.String("status", err.Error()))
		} else {
			zap.L().Info(module.Name(), zap.String("status", "started"))
		}

		defer func(module Module) {
			shutdownTimeout, shutdownCancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
			defer shutdownCancel()

			if err := module.Shutdown(shutdownTimeout); err != nil {
				zap.L().Error(module.Name(), zap.String("status", err.Error()))
			} else {
				zap.L().Info(module.Name(), zap.String("status", "stopped"))
			}
		}(module)
	}

	zap.L().Info("App", zap.String("status", "started"))

	<-appCtx.Done()

	zap.L().Info("App", zap.String("status", "stopping"))

}
