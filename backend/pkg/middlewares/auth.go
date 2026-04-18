package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/config"
)

func AuthMiddleware(conf *config.Config) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            var tokenString string

            // Coba dari header dulu
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader != "" {
                tokenString = strings.TrimPrefix(authHeader, "Bearer ")
                if tokenString == authHeader {
                    return c.JSON(http.StatusUnauthorized, map[string]string{
                        "message": "invalid authorization format",
                    })
                }
            } else {
                // Fallback ke cookie
                cookie, err := c.Cookie("ws_session")
                if err != nil {
                    return c.JSON(http.StatusUnauthorized, map[string]string{
                        "message": "missing authorization",
                    })
                }
                tokenString = cookie.Value
            }

            // Verify token
            claims, err := config.ValidationToken(tokenString, conf)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "message": err.Error(),
                })
            }

            c.Set("user_id", claims.UserID)
            c.Set("username", claims.Username)
            c.Set("role", claims.Role)
            c.Set("session_id", claims.SessionID)
            

            return next(c)
        }
    }
}

