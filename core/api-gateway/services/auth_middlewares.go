package services

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/core/api-gateway/config"
	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
	"github.com/wafi11/workspaces/core/api-gateway/pkg/validate"
)

func (service *AuthService) AuthMiddleware(conf *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			cookie, err := c.Cookie(pkg.CookieAccessTokenName)
			if err != nil {
				log.Printf("errr : %s", err.Error())
				return validate.HandleAuthError(c, err)
			}
			tokenString = cookie.Value

			// Verify token
			claims, err := service.client.ValidateToken(c.Request().Context(), &v1.ValidateTokenRequest{
				Token: tokenString,
			})
			if err != nil {
				return validate.HandleAuthError(c, err)
			}

			c.Set("user_id", claims.UserId)
			c.Set("role", claims.Role)
			c.Set("session_id", claims.SessionId)

			return next(c)
		}
	}
}
