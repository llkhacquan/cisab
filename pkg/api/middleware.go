package api

import (
	"net/http"
	"time"

	"github.com/llkhacquan/cisab/pkg/dbctx"
	"github.com/llkhacquan/cisab/pkg/utils/logger"
	"gorm.io/gorm"
)

// LoggingMiddleware logs information about each HTTP request
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Info("request started", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr, "user_agent", r.UserAgent())
			next.ServeHTTP(w, r)
			log.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
		})
	}
}

// CorsMiddleware adds CORS headers to the response.
// It allows requests from any origin and supports credentials (for development purposes).
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Origin header from the request
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")

		// Allow common headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")

		// Set max age for preflight cache (1 hour)
		w.Header().Set("Access-Control-Max-Age", "3600")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func DBTransactionMiddleware(l *logger.Logger, db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a response recorder to capture the status code
			rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			tx := db.Begin()
			r = r.WithContext(dbctx.Set(r.Context(), tx))
			next.ServeHTTP(rec, r)
			if rec.statusCode != http.StatusOK {
				tx.Rollback()
			} else {
				if err := tx.Commit().Error; err != nil {
					http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
					l.Error("failed to commit transaction", "error", err)
				}
			}
		})
	}
}

// responseRecorder is a custom http.ResponseWriter to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
