package middleware

import (
	"net/http"

	"github.com/v-vovk/health-tracker-api/internal/infra/errors"
	"github.com/v-vovk/health-tracker-api/internal/infra/logger"
	"go.uber.org/zap"
)

// RecoveryMiddleware recovers from panics and logs the error.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error("Recovered from panic", zap.Any("error", err))
				errors.JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
