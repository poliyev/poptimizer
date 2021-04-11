package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r := chi.NewRouter()
	// TODO: добавить другие middleware
	r.Use(s.logger)
	r.Use(middleware.RedirectSlashes)
	r.Get("/{group}/{name}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
		defer cancel()
		res, err := s.repo.ViewJOSN(ctx, domain.TableID{domain.Group(chi.URLParam(r, "group")), domain.Name(chi.URLParam(r, "name"))})
		if err != nil {
			zap.L().Panic(s.Name(), zap.Error(err))
		}
		_, err = w.Write(res)
		if err != nil {
			zap.L().Panic(s.Name(), zap.Error(err))
		}
	})

	// Как писать JSON
	//func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	//buf := &bytes.Buffer{}
	//enc := json.NewEncoder(buf)
	//enc.SetEscapeHTML(true)
	//if err := enc.Encode(v); err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//return
	//}
	//
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
	//w.WriteHeader(status)
	//}
	//w.Write(buf.Bytes())
	//}

	// Как сообщать об ошибке
	// https://golang.org/src/net/http/server.go?s=64501:64553#L2068

	// Добавить настройки кастомного сервера
	srv := http.Server{Addr: s.addr, Handler: r}

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
