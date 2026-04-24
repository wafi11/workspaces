package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/config"
)

func setAuthCookies(c echo.Context, accessToken, refreshToken string) {
	domain := fmt.Sprintf(".%s", "wfdnstore.online")
	c.SetCookie(&http.Cookie{
		Name:     config.CookieAccessTokenName,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     config.CookieRefreshTokenName,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	})
}
