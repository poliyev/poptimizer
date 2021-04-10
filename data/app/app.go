package app

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"syscall"
	"time"
)

//type app struct {
//	repo *adapters.Repo
//}
//
//func (a *app) Run(ctx context.Context) {

//
//	bus := Bus{repo: a.repo}
//
//	steps := []interface{}{
//		// Источники команд
//		&domain.CheckTradingDay{},
//		// Правила
//
//		// Потребители сообщений
//	}
//	for _, step := range steps {
//		bus.register(step)
//	}
//
//	bus.Run(ctx)
//}
//
//func (a app) GetJson(ctx context.Context) ([]byte, error) {
//	return a.repo.ViewJOSN(ctx, domain.TableID{"trading_dates", "trading_dates"})
//}

type Config struct {
	StartTimeout    time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	ServerAddr      string
	MongoURI        string
	MongoDB         string
}

type Module interface {
	Name() string
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type app struct {
	startTimeout    time.Duration
	shutdownTimeout time.Duration
	modules         []Module
}

func (a *app) Run() {
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

func appTerminationCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 2)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	return ctx
}

func NewServer(cfg Config) *app {
	iss := adapters.NewISSClient()
	factory := domain.NewMainFactory(iss)
	repo := adapters.NewRepo(cfg.MongoURI, cfg.MongoDB, factory)

	modules := []Module{
		&Logger{},
		repo,
		&Server{addr: cfg.ServerAddr, requestTimeout: cfg.RequestTimeout, repo: repo},
	}

	return &app{
		startTimeout:    cfg.StartTimeout,
		shutdownTimeout: cfg.ShutdownTimeout,
		modules:         modules,
	}
}
