package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
	"github.com/wafi11/workspaces/core/api-gateway/pkg/validate"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(as *AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) HandleRegister(c echo.Context) error {

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return pkg.Error(c, http.StatusBadRequest, "Invalid Body Request", nil)
	}

	// 2. Panggil gRPC service (context otomatis dibawa dari Echo)
	res, err := h.authService.Register(c.Request().Context(), &v1.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return validate.HandleAuthError(c, err)
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
		return validate.HandleAuthError(c, err)
	}

	pkg.SetAuthCookies(c, res.AccessToken, res.RefreshToken, true)

	// 3. Kembalikan response sukses
	return pkg.Success(c, http.StatusCreated, "Successfully Login", nil)
}
