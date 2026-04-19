package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/config"
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

	c.SetCookie(&http.Cookie{
		Name:     config.CookieAccessTokenName,
		Value:    resp.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Domain:   fmt.Sprintf(".%s","wfdnstore.online"),
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     config.CookieRefreshTokenName,
		Value:    resp.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Domain:   fmt.Sprintf(".%s","wfdnstore.online"),
		SameSite: http.SameSiteLaxMode,
	})

	return response.Success(c, http.StatusOK, "login success", nil)
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
		log.Printf("failed to refresh token %s",err.Error())
		return response.Error(c, http.StatusUnauthorized, "refresh token failed", err)
	}

	return response.Success(c, http.StatusOK, "token refreshed", resp)
}

// Handler
func (h *Handler) Validate(c echo.Context) error {
    cookie, err := c.Cookie(config.CookieAccessTokenName)
    if err != nil {
        return c.NoContent(http.StatusUnauthorized)
    }

    valid, err := h.services.Validate(c.Request().Context(), cookie.Value)
    if err != nil || !valid {
        return c.NoContent(http.StatusUnauthorized)
    }

	

    return c.NoContent(http.StatusOK)
}

func (h *Handler)  CreatePAT(c echo.Context) error{
	userId := c.Get("user_id").(string)
	var req struct {
		Name string `json:"name"`
		ExpiresAt  time.Time `json:"expires_at"`
	}

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}

	data,err := h.services.CreatePAT(c.Request().Context(),&authservices.CreatePATRequest{
		Name: req.Name,
		ExpiresAt: &req.ExpiresAt,
		UserId: userId,
	})

	if err != nil {
		return response.Error(c,http.StatusInternalServerError,err.Error(),nil)
	}

	return response.Success(c,http.StatusCreated,"Successfully Create Personal Access Token",data)
}



func (h *Handler) DeletePAT(c echo.Context) error {
	userId := c.Get("user_id").(string)
	patId := c.Param("pat_id")

	err := h.services.DeletePAT(c.Request().Context(),patId,userId)

	if err != nil {
		return response.Error(c,http.StatusInternalServerError,err.Error(),nil)
	}

	return response.Success(c,http.StatusOK,"Successfully Delete PAT",nil)
}


func (h *Handler)  GetAllPAT(c echo.Context) error {
	userId := c.Get("user_id").(string)

	data,err := h.services.GetAllPAT(c.Request().Context(),userId)
	if err != nil {
		return response.Error(c,http.StatusInternalServerError,"Failed to get pat",nil)
	}	

	return response.Success(c,http.StatusOK,"Successfully Get All PAT",data)
}
