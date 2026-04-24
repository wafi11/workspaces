package workspace

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	workspaceservice "github.com/wafi11/workspaces/services/workspaces-service"
)

func (h *Handler) AddCollaborators(c echo.Context) error {
	workspaceId := c.Param("workspace_id")

	var req struct {
		Role  string `json:"role"`
		Email string `json:"email"`
	}

	user_id := c.Get("user_id").(string)
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}

	data, err := h.svc.AddCollaborators(c.Request().Context(), workspaceservice.WorkspaceCollaborator{
		Email:       req.Email,
		Role:        req.Role,
		InvitedBy:   user_id,
		WorkspaceId: workspaceId,
	})

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create workspaces", err)
	}

	return response.Success(c, http.StatusCreated, "Successfully Invite Collaboborators", data)
}

func (h *Handler) GetWorkspaceCollaboration(c echo.Context) error {
	user_id := c.Get("user_id").(string)

	data, err := h.svc.GetCollaboratedWorkspaces(c.Request().Context(), user_id)

	if err != nil {
		log.Printf("failed to find collaborations teams : %s", err.Error())
		return response.Error(c, http.StatusBadRequest, "Failed to get workspace collaboration", nil)
	}

	return response.Success(c, http.StatusOK, "Successfully retreived workspaces", data)
}

func (h *Handler) AcceptInvite(c echo.Context) error {
	user_id := c.Get("user_id").(string)
	var req struct {
		NotificationID string `json:"notification_id"`
		Types          string `json:"types"`
	}

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err)
	}
	log.Printf("[handler] body: notificationID=%s type=%s", req.NotificationID, req.Types)

	err := h.svc.AcceptOrDeniedInvitationCollborator(c.Request().Context(), req.Types, req.NotificationID, user_id)
	if err != nil {
		return response.Error(c, http.StatusForbidden, "Not a collaborator", err)
	}

	return response.Success(c, http.StatusOK, "Successfuly Accept Notifications", nil)
}
