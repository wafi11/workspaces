package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/response"
)

func AuthMiddleware(conf *config.Config) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            var tokenString string

            cookie, err := c.Cookie(config.CookieAccessTokenName)
            if err != nil {
                return response.Error(c,http.StatusUnauthorized, "Unauthorized",nil)
            }
            tokenString = cookie.Value
        
            // Verify token
            claims, err := config.ValidationToken(tokenString, conf)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "message": err.Error(),
                })
            }

            c.Set("user_id", claims.UserID)
           
            c.Set("role", claims.Role)
            c.Set("session_id", claims.SessionID)
            

            return next(c)
        }
    }
}

