package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
)

func (h *Handler) GetQuota(c echo.Context) error {
	userId := c.Get("user_id").(string)

	if userId == "" {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	data, err := h.svc.GetUserQuota(c.Request().Context(), userId)

	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	return response.Success(c, http.StatusOK, "user quota retreived successfully", data)
}
