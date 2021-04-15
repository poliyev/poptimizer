package ports

import (
	"context"
	"errors"
	"net/http"
	"poptimizer/data/adapters"
	"time"

	"go.uber.org/zap"
)

// Server реализует интерфейс модуля приложения поверх http.Server.
type Server struct {
	http.Server
}

// NewServer тонкая обертка http.Server с интерфейсом модуля приложения.
func NewServer(addr string, requestTimeouts time.Duration, jsonViewer adapters.JSONViewer) *Server {
	srv := Server{
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
func (s *Server) Name() string {
	return "Server"
}

// Start - запускает сервер в отдельной горутине.
func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			zap.L().Panic(s.Name(), zap.Error(err))
		}
	}()

	return nil
}
