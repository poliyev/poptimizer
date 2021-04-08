package main

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"poptimizer/data/app"
	"sync"

	chi "github.com/go-chi/chi/v5"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Printf("\n[WARN] interrupt signal")
		cancel()
	}()
	q := app.App{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		q.Run(ctx)
	}()

	r := chi.NewRouter()
	// Посмотреть и добавить другие middleware
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, _ := q.GetJson(ctx)
		w.Write(res)
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

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}

	wg.Wait()

}
