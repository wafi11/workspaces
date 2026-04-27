package pkg

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	CookieAccessTokenName  string = "workspace_token"
	CookieRefreshTokenName string = "ws_refresh_token"
)

func SetAuthCookies(c echo.Context, accessToken, refreshToken string) {
	domain := fmt.Sprintf(".%s", "wfdnstore.online")
	c.SetCookie(&http.Cookie{
		Name:     CookieAccessTokenName,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     CookieRefreshTokenName,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	})
}
