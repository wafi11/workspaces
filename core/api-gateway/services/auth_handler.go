package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(as *AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) HandleRegister(c echo.Context) error {
	req := new(v1.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return pkg.Error(c, http.StatusBadRequest, "Invalid Body Request", nil)
	}

	// 2. Panggil gRPC service (context otomatis dibawa dari Echo)
	res, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		if err.Error() != "EmailAlreadyExists" {
			return pkg.Error(c, 400, res.Message, err)
		} else if err.Error() == "Invalid Credentials" {
			return pkg.Error(c, 400, "Invalid Credentials", err)
		} else {
			return pkg.Error(c, http.StatusInternalServerError, "Internal Server Error", nil)
		}
	}

	return pkg.Success(c, http.StatusCreated, res.Message, nil)

}

func (h *AuthHandler) HandleLogin(c echo.Context) error {
	// 1. Bind input dari JSON body
	req := new(v1.LoginRequest)
	if err := c.Bind(req); err != nil {
		return pkg.Error(c, http.StatusBadRequest, "Invalid Body Request", nil)
	}

	// 2. Panggil gRPC service (context otomatis dibawa dari Echo)
	res, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "EmailAlreadyExists" {
			return pkg.Error(c, 10, "Email Already Exists", nil)
		} else if err.Error() == "Invalid Credentials" {
			return pkg.Error(c, 400, "Invalid Credentials", nil)
		} else {
			return pkg.Error(c, http.StatusInternalServerError, "Internal Server Error", nil)
		}
	}

	pkg.SetAuthCookies(c, res.AccessToken, res.RefreshToken)

	// 3. Kembalikan response sukses
	return pkg.Success(c, http.StatusCreated, "Successfully Login", nil)
}
