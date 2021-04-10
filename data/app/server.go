package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	addr string
	srv  *http.Server
}

func (s *Server) Name() string {
	return "Server"
}

func (s *Server) Start(ctx context.Context) error {
	r := chi.NewRouter()
	// Посмотреть и добавить другие middleware
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// res, _ := q.GetJson(ctx)
		w.Write([]byte("Hi!"))
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
