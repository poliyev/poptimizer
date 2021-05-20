package ports

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Server реализует интерфейс модуля приложения поверх http.Server.
type Server struct {
	http.Server
}

// NewServer тонкая обертка http.Server с интерфейсом модуля приложения.
func NewServer(addr string, requestTimeouts time.Duration, jsonViewer JSONViewer) *Server {
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

// Start - запускает сервер в отдельной горутине.
func (s *Server) Start(_ context.Context) error {
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			zap.L().Panic("Server", zap.Error(err))
		}
	}()

	return nil
}
