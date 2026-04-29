package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	CookieAccessTokenName  string = "workspace_token"
	CookieRefreshTokenName string = "ws_refresh_token"
)

func SetAuthCookies(c echo.Context, accessToken, refreshToken string, isDev bool) {
	domain := ".wfdnstore.online"
	if isDev {
		domain = "localhost"
	}

	c.SetCookie(&http.Cookie{
		Name:     CookieAccessTokenName,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   !isDev, // false kalau dev
		Path:     "/",
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     CookieRefreshTokenName,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   !isDev,
		Path:     "/",
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	})
}
