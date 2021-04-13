package ports

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"poptimizer/data/adapters"
	"time"
)

// server реализует интерфейс модуля приложения поверх http.Server.
type server struct {
	http.Server
}

// NewServer создает http.Server с интерфейсом модуля приложения.
func NewServer(addr string, requestTimeouts time.Duration, jsonViewer adapters.JSONViewer) *server {
	srv := server{
		Server: http.Server{
			Addr:         addr,
			ReadTimeout:  requestTimeouts,
			WriteTimeout: requestTimeouts,
		},
	}
	srv.Handler = newTableMux(requestTimeouts, jsonViewer)
	return &srv
}

// Name - модуль Server.
func (s *server) Name() string {
	return "Server"
}

// Start - запускает сервер в отдельной горутине.
func (s *server) Start(ctx context.Context) error {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Panic(s.Name(), zap.Error(err))
		}
	}()

	return nil
}
