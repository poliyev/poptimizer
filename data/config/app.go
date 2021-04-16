package config

import (
	"context"
	"os"
	"os/signal"
	"poptimizer/data/adapters"
	"poptimizer/data/app"
	"poptimizer/data/domain"
	"poptimizer/data/ports"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// App - обеспечивает запуск и остановку приложения.
//
// Приложение состоит из отдельных модулей, которые последовательно запускаются в начале (обычно начиная с
// инфраструктуры и заканчивая модулями взаимодействующими с пользователем) и останавливаются в обратном порядке в
// конце. Остановка осуществляется с помощью системных сигналов.
type App struct {
	startTimeout    time.Duration
	shutdownTimeout time.Duration
	stop            chan os.Signal
	modules         []Module
}

// NewApp - создает приложение на основе конфигурации.
func NewApp(cfg *Config) *App {
	iss := adapters.NewISSClient(cfg.ISSMaxCons)
	factory := domain.NewMainFactory(iss)
	repo := adapters.NewRepo(cfg.MongoURI, cfg.MongoDB, factory)
	bus := app.NewBus(repo, cfg.EventBusTimeouts)

	modules := []Module{
		adapters.NewLogger(),
		repo,
		bus,
		ports.NewServer(cfg.ServerAddr, cfg.ServerTimeouts, repo),
	}

	return &App{
		startTimeout:    cfg.StartTimeout,
		shutdownTimeout: cfg.ShutdownTimeout,
		stop:            make(chan os.Signal, 1),
		modules:         modules,
	}
}

// Run - запускает модули приложения, блокируется на получении системных сигналов SIGINT или SIGTERM и осуществляет
// завершение работы модулей после их поступления.
func (a *App) Run() {
	a.startModules()
	a.terminated()
	a.shutdownModules()
}

func (a *App) startModules() {
	startCtx, startCancel := context.WithTimeout(context.Background(), a.startTimeout)
	defer startCancel()

	for _, module := range a.modules {
		if err := module.Start(startCtx); err != nil {
			zap.L().Panic(module.Name(), zap.String("status", err.Error()))
		}

		zap.L().Info(module.Name(), zap.String("status", "started"))
	}

	zap.L().Info("App", zap.String("status", "started"))
}

func (a *App) terminated() {
	signal.Notify(a.stop, syscall.SIGINT, syscall.SIGTERM)
	<-a.stop
	zap.L().Info("App", zap.String("status", "stopping"))
}

func (a *App) shutdownModules() {
	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	modules := a.modules
	for n := range modules {
		module := modules[len(modules)-1-n]

		if err := module.Shutdown(ctx); err != nil {
			zap.L().Warn(module.Name(), zap.String("status", err.Error()))
		} else {
			zap.L().Info(module.Name(), zap.String("status", "stopped"))
		}
	}

	zap.L().Info("App", zap.String("status", "stopped"))
}
