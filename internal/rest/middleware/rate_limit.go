package middleware

import (
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"net/http"
	"time"
)

type RateLimitMiddleware struct {
	renderer *renderer.Renderer
}

func NewRateLimitMiddleware(renderer *renderer.Renderer) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		renderer: renderer,
	}
}

func (m *RateLimitMiddleware) Limit(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return httprate.LimitBy(
		requestLimit,
		windowLength,
		clientIPKey,
		m.renderPage(),
	)
}

func (m *RateLimitMiddleware) LimitByEndpoint(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return httprate.LimitBy(
		requestLimit,
		windowLength,
		httprate.JoinKeys(clientIPKey, httprate.KeyByEndpoint),
		m.renderPage(),
	)
}

func (m *RateLimitMiddleware) renderPage() httprate.Option {
	return httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
		_ = m.renderer.RenderTooManyRequests(w, r)
	})
}

func clientIPKey(r *http.Request) (string, error) {
	return "limit:" + httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
}
