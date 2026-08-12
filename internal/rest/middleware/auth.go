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

type AuthMiddleware struct {
	userSvc  *user.Service
	renderer *renderer.Renderer
	sessCfg  config.Session
	log      *logger.Logger
}

func NewAuthMiddleware(
	userSvc *user.Service,
	renderer *renderer.Renderer,
	sessCfg config.Session,
	log *logger.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		userSvc:  userSvc,
		renderer: renderer,
		sessCfg:  sessCfg,
		log:      log,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	Check whether user unauthorized
		cookieName := rest.GetSessionCookieName(m.sessCfg.Name)
		cookie, err := r.Cookie(cookieName)

		//	Unauthorized -> redirect to sign up page
		if err != nil {
			m.renderer.Redirect(w, r, "/auth/sign-up")
			return
		}

		//	Check if session ID is valid
		u, err := m.userSvc.GetUserBySID(cookie.Value)

		//	Expired -> redirect to sign up page
		if err != nil && errors.Is(err, user.ErrSessionExpired) {
			m.renderer.Redirect(w, r, "/auth/sign-up")
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

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
