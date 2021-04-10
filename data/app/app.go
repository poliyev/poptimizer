package app

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//type app struct {
//	repo *adapters.Repo
//}
//
//func (a *app) Run(ctx context.Context) {
//	if a.repo == nil {
//		iss := adapters.NewISSClient()
//		factory := domain.NewMainFactory(iss)
//		// TODO: закрывать РЕПО
//		a.repo = adapters.NewRepo(ctx, MongoURI, mongoDB, factory)
//	}
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

	appCtx := appTerminationCtx()

	startCtx, startCancel := context.WithTimeout(appCtx, a.startTimeout)
	defer startCancel()

	shutdownTimeout, shutdownCancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer shutdownCancel()

	for _, module := range a.modules {
		if err := module.Start(startCtx); err != nil {
			zap.L().Panic(module.Name(), zap.String("status", err.Error()))
		} else {
			zap.L().Info(module.Name(), zap.String("status", "started"))
		}

		defer func(module Module) {
			if err := module.Shutdown(shutdownTimeout); err != nil {
				zap.L().Error(module.Name(), zap.String("status", err.Error()))
			} else {
				zap.L().Info(module.Name(), zap.String("status", "stopped"))
			}
		}(module)
	}

	<-appCtx.Done()
}

func appTerminationCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 2)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stop
		zap.L().Info("App", zap.String("status", sig.String()))
		cancel()
	}()

	return ctx
}

func NewServer(cfg Config) *app {
	modules := []Module{
		Logger{},
		&Server{addr: cfg.ServerAddr}}

	return &app{
		startTimeout:    cfg.StartTimeout,
		shutdownTimeout: cfg.ShutdownTimeout,
		modules:         modules,
	}
}
