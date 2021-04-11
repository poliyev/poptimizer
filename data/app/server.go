package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"time"
)

type Server struct {
	addr           string
	requestTimeout time.Duration
	srv            *http.Server
	repo           *adapters.Repo
}

func (s *Server) Name() string {
	return "Server"
}

func (s *Server) logger(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			zap.L().Info(
				s.Name(),
				zap.String("request", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
				zap.Duration("time", time.Since(start)))
		}()

		handler.ServeHTTP(ww, r)

	})

}

func (s *Server) Start(ctx context.Context) error {
	router := chi.NewRouter()
	// TODO: добавить другие middleware
	router.Use(s.logger)
	router.Use(middleware.RedirectSlashes)
	router.Get("/{group}/{name}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
		defer cancel()
		res, err := s.repo.ViewJOSN(ctx, domain.TableID{domain.Group(chi.URLParam(r, "group")), domain.Name(chi.URLParam(r, "name"))})
		if err == mongo.ErrNoDocuments {
			router.NotFoundHandler()(w, r)
			return
		} else if err != nil {
			// Как сообщать об ошибке
			// https://golang.org/src/net/http/server.go?s=64501:64553#L2068
			zap.L().Panic(s.Name(), zap.Error(err))
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write(res)
		if err != nil {
			zap.L().Panic(s.Name(), zap.Error(err))
		}
	})

	// Добавить настройки кастомного сервера
	srv := http.Server{Addr: s.addr, Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Panic(s.Name(), zap.String("status", err.Error()))
		}
	}()

	s.srv = &srv

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
