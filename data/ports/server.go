package ports

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"poptimizer/data/adapters"
	"time"
)

type server struct {
	http.Server
}

func NewServer(addr string, requestTimeout time.Duration, jsonViewer adapters.JSONViewer) *server {
	srv := server{
		Server: http.Server{
			Addr:         addr,
			ReadTimeout:  requestTimeout,
			WriteTimeout: requestTimeout,
		},
	}
	srv.Handler = NewTableMux(srv.Name(), requestTimeout, jsonViewer)
	return &srv
}

func (s *server) Name() string {
	return "Server"
}

func (s *server) Start(ctx context.Context) error {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Panic(s.Name(), zap.Error(err))
		}
	}()

	return nil
}
