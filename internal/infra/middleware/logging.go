package middleware

import (
	"net/http"
	"time"

	"github.com/v-vovk/health-tracker-api/internal/infra/logger"
	"go.uber.org/zap"
)

// RequestLogger logs details about each HTTP request.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		logger.Log.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", duration),
		)
	})
}
