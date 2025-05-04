package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type loggingResponseWriter struct {
	ResponseWriter http.ResponseWriter
	statusCode     int
	size           int64
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += int64(size)
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.ResponseWriter.Header()
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logOpts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler := slog.NewJSONHandler(os.Stderr, logOpts)
		logger := slog.New(handler)

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		next.ServeHTTP(lrw, r)

		logAttrs := []slog.Attr{
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			),
			slog.Group("response",
				slog.Int("status", lrw.statusCode),
				slog.Duration("duration", time.Since(start)),
				slog.Int64("size", lrw.size),
			),
		}

		if lrw.statusCode == http.StatusOK {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "API Request", logAttrs...)
		} else {
			logger.LogAttrs(context.Background(), slog.LevelWarn, "API Request Error", logAttrs...)
		}
	})
}
