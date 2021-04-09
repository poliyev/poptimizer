package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

func RunServer() *http.Server {
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
	srv := http.Server{Addr: ":3000", Handler: r}

	go func() {
		zap.L().Info("Сервер запущен")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			zap.L().Panic("Ошибка при работе сервера", zap.Error(err))
		}
	}()

	return &srv
}

func StopServer(srv *http.Server) {
	if err := srv.Shutdown(context.Background()); err != nil {
		zap.L().Error("Ошибка остановки сервера", zap.Error(err))
	} else {
		zap.L().Info("Сервер остановлен")
	}
}
