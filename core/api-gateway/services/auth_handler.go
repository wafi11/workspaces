package services

import (
	"log"
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
		Name     string `json:"name"`
	}

	if err := c.Bind(&req); err != nil {
		return pkg.Error(c, http.StatusBadRequest, "Invalid Body Request", nil)
	}

	// 2. Panggil gRPC service (context otomatis dibawa dari Echo)
	res, err := h.authService.Register(c.Request().Context(), &v1.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
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

func (h *AuthHandler) HandleRefreshToken(c echo.Context) error {

	var tokenString string

	cookie, err := c.Cookie(pkg.CookieRefreshTokenName)
	if err != nil {
		log.Printf("errr : %s", err.Error())
		return validate.HandleAuthError(c, err)
	}
	tokenString = cookie.Value

	res, err := h.authService.RefreshToken(c.Request().Context(), &v1.RefreshTokenRequest{
		RefreshToken: tokenString,
	})

	if err != nil {
		log.Printf("[Refresh Token] failed to refresh token : %s", err.Error())
		return pkg.Error(c, http.StatusInternalServerError, "failed to refresh token", nil)
	}

	pkg.SetAuthCookies(c, res.AccessToken, tokenString, true)

	return pkg.Success(c, http.StatusOK, "Successfully Refresh Token", nil)
}

func (h *AuthHandler) GetOAuthURLGithub(c echo.Context) error {

	res, err := h.authService.GetOAuthURL(c.Request().Context(), &v1.GetOAuthURLRequest{
		Provider: "github",
	})

	if err != nil {
		return pkg.Error(c, http.StatusInternalServerError, "internal server error", nil)
	}

	return pkg.Success(c, http.StatusOK, "Successfully Get Oauth Url", res.Url)

}

func (h *AuthHandler) ConnectOAuthGithub(c echo.Context) error {

	userId := c.Get("user_id").(string)
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	_, err := h.authService.ConnectOAuth(c.Request().Context(), &v1.ConnectOAuthRequest{
		UserId:   userId,
		Provider: "github",
		Code:     code,
		State:    state,
	})

	if err != nil {
		return pkg.Error(c, http.StatusInternalServerError, "internal server error", nil)
	}

	return pkg.Success(c, http.StatusOK, "Successfully Connect Oauth", nil)
}

func (h *AuthHandler) OAuthLoginGithub(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	res, err := h.authService.OAuthLogin(c.Request().Context(), &v1.OAuthCallbackRequest{
		Provider: "github",
		Code:     code,
		State:    state,
	})

	if err != nil {
		return pkg.Error(c, http.StatusInternalServerError, "internal server error", nil)
	}

	pkg.SetAuthCookies(c, res.AccessToken, res.RefreshToken, true)

	return pkg.Success(c, http.StatusOK, "Successfully Login", nil)
}
