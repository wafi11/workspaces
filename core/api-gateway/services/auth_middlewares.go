package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/core/api-gateway/config"
	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
	"github.com/wafi11/workspaces/core/api-gateway/pkg/validate"
)

func (service *AuthService) AuthMiddleware(conf *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(pkg.CookieAccessTokenName)
			if err != nil {
				// http.ErrNoCookie adalah error standard dari net/http
				return pkg.Error(c, http.StatusUnauthorized, "unauthorized", nil)
			}

			if cookie.Value == "" {
				return pkg.Error(c, http.StatusUnauthorized, "unauthorized", nil)
			}

			claims, err := service.client.ValidateToken(c.Request().Context(), &v1.ValidateTokenRequest{
				Token: cookie.Value,
			})
			if err != nil {
				return validate.HandleAuthError(c, err) // ini udah pasti gRPC error
			}

			c.Set("user_id", claims.UserId)
			c.Set("role", claims.Role)
			c.Set("session_id", claims.SessionId)

			return next(c)
		}
	}
}
