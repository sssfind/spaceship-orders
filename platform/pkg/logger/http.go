package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// HTTP логирует входящие запросы через Zap (stdout JSON и OTLP).
// /health и /metrics пропускаются, чтобы probes не забивали Victoria Logs.
func HTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" || r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(sw, r)

		fields := []zap.Field{
			zap.String("http.method", r.Method),
			zap.String("http.path", r.URL.Path),
			zap.Int("http.status", sw.status),
			zap.Duration("duration", time.Since(start)),
		}

		ctx := r.Context()
		msg := "request completed"
		switch {
		case sw.status >= 500:
			Error(ctx, msg, fields...)
		case sw.status >= 400:
			Warn(ctx, msg, fields...)
		default:
			Info(ctx, msg, fields...)
		}
	})
}
