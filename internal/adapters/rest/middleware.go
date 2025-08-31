package rest

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/Util787/task-processor/internal/common"
)

const maxExpectedDurationMs = 2000

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: 0}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func basicMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()

		ctx := context.WithValue(r.Context(), common.ContextKey("request_id"), requestID)
		r = r.WithContext(ctx)

		logger := log.With(slog.String("request_id", requestID))

		logger.Info("Request received", slog.String("ip", r.RemoteAddr), slog.String("user_agent", r.UserAgent()), slog.String("path", r.URL.Path))

		start := time.Now()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		durationMs := time.Since(start).Milliseconds()
		logger.Debug("Request finished", slog.Int("status_code", rw.statusCode), slog.Int64("duration_ms", durationMs))
		if durationMs > maxExpectedDurationMs {
			logger.Warn("Operation is taking more time than expected", slog.Int("expected_duration(ms)", maxExpectedDurationMs), slog.Int64("actual_duration(ms)", durationMs))
		}
	})
}

func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().Unix(), rand.Intn(100000000))
}

func validateMethodMiddleware(allowedMethod string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
