package api

import (
	"net/http"
	"time"

	"github.com/llkhacquan/knovel-assignment/pkg/dbctx"
	"github.com/llkhacquan/knovel-assignment/pkg/utils/logger"
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

// CorsMiddleware adds CORS headers to the response
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware is a simple auth middleware that checks for an API key
// For demonstration purposes only - not secure for production
func AuthMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// todo: implement actual authentication logic
			// now we just assume the request is authenticated
			next.ServeHTTP(w, r)
		})
	}
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
