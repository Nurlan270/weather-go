package middleware

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/entity"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest"
	"github.com/Nurlan270/weather-go/internal/user"
)

type mockUserRepository struct{}

func (m mockUserRepository) GetUserBySID(sid string) (*entity.User, error) {
	return nil, sql.ErrNoRows
}

func (m mockUserRepository) GetUserByLogin(login string) (*entity.User, error) {
	panic("implement me")
}

func (m mockUserRepository) CreateUser(login, password string) (*entity.User, error) {
	panic("implement me")
}

func (m mockUserRepository) CreateSession(
	uuid string,
	userID int64,
	expiresAt time.Time,
) (*entity.Session, error) {
	panic("implement me")
}

func TestAuthMiddleware(t *testing.T) {
	//	Arrange
	sessConf := config.Session{
		Name:      "test",
		ExpiresIn: 5 * time.Minute,
	}

	//	Tests
	t.Run("it redirects unauthorized user to sign up page", func(t *testing.T) {
		//	Arrange
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		mw := NewAuthMiddleware(
			&user.Service{},
			renderer.New(nil),
			sessConf,
			logger.New(logger.EnvLocal),
		)

		//	Imitate that this is route that requires authorized user
		nextCalled := false
		next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		})

		//	Act
		mw.RequireAuth(next).ServeHTTP(rr, req)

		require.False(t, nextCalled, "Should not call next handler")
		require.Equal(t, http.StatusSeeOther, rr.Code)

		expected := "/auth/sign-up"
		got := rr.Header().Get("Location")
		require.Equal(t, expected, got, "Should redirect user to sign up page")
	})

	t.Run("it redirects user with expired session cookie to sign up page", func(t *testing.T) {
		//	Arrange
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  rest.GetSessionCookieName(sessConf.Name),
			Value: "expired-session-id",
		})

		rr := httptest.NewRecorder()

		mw := NewAuthMiddleware(
			user.NewService(mockUserRepository{}, sessConf),
			renderer.New(nil),
			sessConf,
			logger.New(logger.EnvLocal),
		)

		//	Imitate that this is handler that requires authorized user
		nextCalled := false
		next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		})

		//	Act
		mw.RequireAuth(next).ServeHTTP(rr, req)

		require.False(t, nextCalled, "Should not call next handler")
		require.Equal(t, http.StatusSeeOther, rr.Code)

		expected := "/auth/sign-up"
		got := rr.Header().Get("Location")
		require.Equal(t, expected, got, "Should redirect user to sign up page")
	})
}
