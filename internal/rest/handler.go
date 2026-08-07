package rest

import (
	"github.com/iancoleman/strcase"
	"net/http"
	"strings"
	"time"
)

// GetSessionCookieName builds session cookie's name
func GetSessionCookieName(sessName string) string {
	s := strcase.ToSnake(strings.ToLower(sessName))

	return s + "_session"
}

func SetCookie(w http.ResponseWriter, name, value string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
