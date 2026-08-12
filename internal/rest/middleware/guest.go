package middleware

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest"
	"github.com/Nurlan270/weather-go/internal/user"
)

type GuestMiddleware struct {
	userSvc  *user.Service
	renderer *renderer.Renderer
	sessCfg  config.Session
	log      *logger.Logger
}

func NewGuestMiddleware(
	userSvc *user.Service,
	renderer *renderer.Renderer,
	sessCfg config.Session,
	log *logger.Logger,
) *GuestMiddleware {
	return &GuestMiddleware{
		userSvc:  userSvc,
		renderer: renderer,
		sessCfg:  sessCfg,
		log:      log,
	}
}

func (m *GuestMiddleware) RequireGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieName := rest.GetSessionCookieName(m.sessCfg.Name)
		cookie, err := r.Cookie(cookieName)

		//	Unauthorized -> continue
		if err != nil {
			next.ServeHTTP(w, r)

			return
		}

		//	Check if session ID is valid
		u, err := m.userSvc.GetUserBySID(cookie.Value)

		//	Session expired -> continue
		if err != nil && errors.Is(err, user.ErrSessionExpired) {
			next.ServeHTTP(w, r)

			return
		}

		if err != nil {
			m.log.Error("could not get user by session ID", zap.Error(err))

			if err = m.renderer.RenderServerError(w, u.Login); err != nil {
				m.log.ErrorRenderPage(err)
			}

			return
		}

		//	Set user to ctx, so it's accessible from handlers
		ctx := user.NewContext(r.Context(), u)

		//	Authorized -> redirect to home page
		m.renderer.Redirect(w, r.WithContext(ctx), "/")
	})
}
