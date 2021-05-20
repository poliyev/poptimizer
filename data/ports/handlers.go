package ports

import (
	"context"
	"errors"
	"net/http"
	"poptimizer/data/domain"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// JSONViewer обеспечивает просмотр данных таблиц в формате ExtendedJSON.
type JSONViewer interface {
	ViewJSON(ctx context.Context, id domain.ID) ([]byte, error)
}

type tableHandler struct {
	viewer JSONViewer
}

func (t tableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := domain.NewID(chi.URLParam(r, "group"), chi.URLParam(r, "name"))
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

func newTableMux(requestTimeouts time.Duration, viewer JSONViewer) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Timeout(requestTimeouts))
	router.Use(zapLoggingMiddleware)
	router.Use(middleware.RedirectSlashes)
	router.Method(http.MethodGet, "/{group}/{name}", tableHandler{viewer})

	return router
}
