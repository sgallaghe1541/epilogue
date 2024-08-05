package auth

import "net/http"

const (
	JWTCookie     = "token"
	RefreshCookie = "refresh"
)

func NewJWTCookie(name, value string) string {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie.String()
}

func NewRefreshCookie(name, value string) string {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   2592000,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie.String()
}
