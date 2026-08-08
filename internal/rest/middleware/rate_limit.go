package middleware

import (
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"net/http"
	"time"
)

type Limiter struct {
	renderer *renderer.Renderer
}

func NewRateLimitMiddleware(renderer *renderer.Renderer) *Limiter {
	return &Limiter{
		renderer: renderer,
	}
}

func (l *Limiter) Limit(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return httprate.LimitBy(
		requestLimit,
		windowLength,
		clientIPKey,
		l.renderPage(),
	)
}

func (l *Limiter) LimitByEndpoint(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return httprate.LimitBy(
		requestLimit,
		windowLength,
		httprate.JoinKeys(clientIPKey, httprate.KeyByEndpoint),
		l.renderPage(),
	)
}

func (l *Limiter) renderPage() httprate.Option {
	return httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
		_ = l.renderer.RenderTooManyRequests(w, r)
	})
}

func clientIPKey(r *http.Request) (string, error) {
	return "limit:" + httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
}
