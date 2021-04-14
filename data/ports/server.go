package ports

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
	"poptimizer/data/adapters"
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
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			zap.L().Panic(s.Name(), zap.Error(err))
		}
	}()

	return nil
}
