package ports

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"time"
)

type TableHandler struct {
	serverName string
	viewer     adapters.JSONViewer
}

func (t TableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := domain.NewTableID(chi.URLParam(r, "group"), chi.URLParam(r, "name"))
	res, err := t.viewer.ViewJSON(ctx, id)
	if err == mongo.ErrNoDocuments {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	} else if err != nil {
		// https://golang.org/src/net/http/server.go?s=64501:64553#L2068
		zap.L().Panic(t.serverName, zap.Error(err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(res)
	if err != nil {
		zap.L().Panic(t.serverName, zap.Error(err))
	}
}

func NewTableMux(serverName string, requestTimeout time.Duration, viewer adapters.JSONViewer) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Timeout(requestTimeout))
	router.Use(ZapLoggingMiddleware(serverName))
	router.Use(middleware.RedirectSlashes)
	router.Method(http.MethodGet, "/{group}/{name}", TableHandler{serverName, viewer})

	return router
}
