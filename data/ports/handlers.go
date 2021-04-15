package ports

import (
	"errors"
	"net/http"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type tableHandler struct {
	viewer adapters.JSONViewer
}

func (t tableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := domain.NewTableID(chi.URLParam(r, "group"), chi.URLParam(r, "name"))
	res, err := t.viewer.ViewJSON(ctx, id)

	if errors.Is(err, mongo.ErrNoDocuments) {
		http.NotFound(w, r)
		zap.L().Warn("JSONViewer", zap.Error(err))

		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zap.L().Warn("JSONViewer", zap.Error(err))

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err = w.Write(res); err != nil {
		zap.L().Warn("tableHandler", zap.Error(err))
	}
}

func newTableMux(requestTimeouts time.Duration, viewer adapters.JSONViewer) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Timeout(requestTimeouts))
	router.Use(zapLoggingMiddleware)
	router.Use(middleware.RedirectSlashes)
	router.Method(http.MethodGet, "/{group}/{name}", tableHandler{viewer})

	return router
}
