package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	userservices "github.com/wafi11/workspaces/services/user-service"
)

type Handler struct {
	svc userservices.UserService
}

func NewHandler(svc userservices.UserService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProfile(c echo.Context) error {
	user_id := c.Get("user_id").(string)

	res, err := h.svc.GetProfile(c.Request().Context(), &userservices.GetUserRequest{
		UserId: user_id,
	})
	if err != nil {
		return response.Error(c, http.StatusNotFound, "user not found", err)
	}

	return response.Success(c, http.StatusOK, "success", res.User)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	var req userservices.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}

	req.UserId = c.Get("user_id").(string)

	res, err := h.svc.UpdateUser(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to update user", err)
	}

	return response.Success(c, http.StatusOK, "user updated", res)
}

func (h *Handler) ChangePassword(c echo.Context) error {
	var req userservices.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}
	req.UserId = c.Get("user_id").(string)

	res, err := h.svc.ChangePassword(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to change password", err)
	}

	return response.Success(c, http.StatusOK, "password changed", res)
}

func (h *Handler) GetUserSessions(c echo.Context) error {
	req := userservices.GetUserSessionsRequest{
		UserId: c.Get("user_id").(string),
	}

	res, err := h.svc.GetUserSessions(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to get sessions", err)
	}

	return response.Success(c, http.StatusOK, "success", res)
}

func (h *Handler) RevokeSession(c echo.Context) error {
	req := userservices.RevokeSessionRequest{
		SessionId: c.Param("session_id"),
	}

	res, err := h.svc.RevokeSession(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to revoke session", err)
	}

	return response.Success(c, http.StatusOK, "session revoked", res)
}
