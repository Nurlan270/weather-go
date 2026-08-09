package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/Nurlan270/weather-go/internal/logger"
)

type LoggerMiddleware struct {
	log *logger.Logger
}

func NewLoggerMiddleware(log *logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{log: log}
}

func (m *LoggerMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.log.With(
			zap.String("ip", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.String("request_id", middleware.GetReqID(r.Context())),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Time("timestamp", time.Now().UTC()),
		)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			log.Info("Request information",
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
				zap.Duration("duration", time.Since(t1).Round(time.Millisecond)),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
