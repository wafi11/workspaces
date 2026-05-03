package services

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(svc *UserService) *UserHandler {
	return &UserHandler{
		service: svc,
	}
}
func (h *UserHandler) GetProfile(c echo.Context) error {
	// Aman dari panic kalau middleware ga set value
	userId, ok := c.Get("user_id").(string)
	if !ok || userId == "" {
		return pkg.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	sessionId, ok := c.Get("session_id").(string)
	if !ok || sessionId == "" {
		return pkg.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	profile, err := h.service.GetProfile(c.Request().Context(), &v1.GetProfileRequest{
		Id:        userId,
		SessionId: sessionId,
	})
	if err != nil {
		log.Printf("[profile] error: %s", err.Error())
		return pkg.Error(c, http.StatusInternalServerError, "failed to get profile", nil)
	}

	return pkg.Success(c, http.StatusOK, "Successfully Get Profile", profile.User)
}
func (h *UserHandler) UpdateProfile(c echo.Context) error {

	req := new(v1.UpdateProfileRequest)
	if err := c.Bind(req); err != nil {
		return pkg.Error(c, http.StatusBadRequest, "Invalid Body Request", nil)
	}

	req.UserId = c.Get("user_id").(string)
	res, err := h.service.UpdateProfile(c.Request().Context(), &v1.UpdateProfileRequest{
		AvatarBase64: req.AvatarBase64,
		Name:         req.Name,
		UserId:       req.UserId,
	})

	if err != nil {

		return pkg.Error(c, http.StatusInternalServerError, "Failed to update", nil)
	}

	return pkg.Success(c, http.StatusOK, "Successfully Update Profile", res)
}
