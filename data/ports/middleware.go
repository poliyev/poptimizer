package ports

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func zapLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			zap.L().Info(
				"Request",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
				zap.Duration("time", time.Since(start)))
		}()

		next.ServeHTTP(ww, r)

	}
	return http.HandlerFunc(fn)
}
