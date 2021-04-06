package main

import (
	"context"
	"net/http"
	"sync"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"os"
	"os/signal"
	"poptimizer/data/app"
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
