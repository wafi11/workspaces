package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	authservices "github.com/wafi11/workspaces/services/auth-service"
)

type Handler struct {
	services authservices.IServices
}

func NewHandler(services authservices.IServices) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Login(c echo.Context) error {
	var req authservices.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}

	userAgent := c.Request().Header.Get("User-Agent")
	ipAddress := c.RealIP()

	resp, err := h.services.Login(c.Request().Context(), &req, userAgent, ipAddress)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "login failed", err)
	}

	return response.Success(c, http.StatusOK, "login success", resp)
}

func (h *Handler) Register(c echo.Context) error {
	var req authservices.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}

	resp, err := h.services.Register(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "register failed", err)
	}

	return response.Success(c, http.StatusCreated, "register success", resp)
}

func (h *Handler) Logout(c echo.Context) error {
	var req authservices.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}

	resp, err := h.services.Logout(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "logout failed", err)
	}

	return response.Success(c, http.StatusOK, "logout success", resp)
}

func (h *Handler) RefreshToken(c echo.Context) error {
	var req authservices.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}

	resp, err := h.services.RefreshToken(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "refresh token failed", err)
	}

	return response.Success(c, http.StatusOK, "token refreshed", resp)
}
