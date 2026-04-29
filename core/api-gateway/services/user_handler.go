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

	userId := c.Get("user_id").(string)
	sesssionId := c.Get("session_id").(string)

	profile, err := h.service.GetProfile(c.Request().Context(), &v1.GetProfileRequest{
		Id:        userId,
		SessionId: sesssionId,
	})

	if err != nil {
		log.Printf("[profile] error : %s", err.Error())
		return pkg.Error(c, http.StatusInternalServerError, "failed to get profile", nil)
	}

	return pkg.Success(c, http.StatusOK, "Successfully Get Profile", profile.User)
}
