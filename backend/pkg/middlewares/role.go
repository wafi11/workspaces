package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/response"
)

func AdminMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            role, _ := c.Get("role").(string)
            if role != string(config.RoleAdmin) {
                return response.Error(c, http.StatusForbidden, "Forbidden", nil)
            }
            return next(c)
        }
    }
}