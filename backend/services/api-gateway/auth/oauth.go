package auth

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
)

// GET /auth/github
func (h *Handler) GithubLogin(c echo.Context) error {
	state := h.services.GenerateState()
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	url := h.services.GithubOauthConfig().AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /auth/github/callback
func (h *Handler) GithubCallback(c echo.Context) error {
	// Verify state
	cookie, err := c.Cookie("oauth_state")
	if err != nil || cookie.Value != c.QueryParam("state") {
		return response.Error(c, http.StatusBadRequest, "invalid oauth state", nil)
	}

	cfg := h.services.GithubOauthConfig()
	token, err := cfg.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "failed to exchange token", err)
	}

	resp, err := h.services.LoginWithGithub(c.Request().Context(), token.AccessToken)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "github login failed", err)
	}

	setAuthCookies(c, resp.AccessToken, resp.RefreshToken)
	return c.Redirect(http.StatusTemporaryRedirect, "https://web-platform.wfdnstore.online/dashboard")
}

// GET /auth/google
func (h *Handler) GoogleLogin(c echo.Context) error {
	state := h.services.GenerateState()
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	url := h.services.GoogleOauthConfig().AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /auth/google/callback
func (h *Handler) GoogleCallback(c echo.Context) error {
	cookie, err := c.Cookie("oauth_state")
	if err != nil || cookie.Value != c.QueryParam("state") {
		return response.Error(c, http.StatusBadRequest, "invalid oauth state", nil)
	}

	cfg := h.services.GoogleOauthConfig()
	token, err := cfg.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "failed to exchange token", err)
	}

	resp, err := h.services.LoginWithGoogle(c.Request().Context(), token.AccessToken)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "google login failed", err)
	}

	setAuthCookies(c, resp.AccessToken, resp.RefreshToken)
	return c.Redirect(http.StatusTemporaryRedirect, "https://web-platform.wfdnstore.online")
}
