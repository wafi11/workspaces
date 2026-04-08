package workspace

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	workspaceservice "github.com/wafi11/workspaces/services/workspaces-service"
)

func (h *Handler) CreateAddonService(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		return response.Error(c, http.StatusNotFound, "workspace not found", nil)
	}

	var req workspaceservice.CreateWorkspaceAddon
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", nil)
	}

	req.WorkspaceID = workspaceID

	if err := h.svc.CreateAddonWorkspace(c.Request().Context(), req); err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to create addon", nil)
	}

	return response.Success(c, http.StatusCreated, "addon created successfully", nil)
}

func (h *Handler) GetAddonService(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		return response.Error(c, http.StatusNotFound, "workspace not found", nil)
	}

	data, err := h.svc.GetAddonService(c.Request().Context(), workspaceID)
	if err != nil {
		if errors.Is(err, workspaceservice.ErrAddonNotFound) {
			return response.Error(c, http.StatusNotFound, "no addons found for this workspace", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "failed to retrieve addons", nil)
	}

	return response.Success(c, http.StatusOK, "addons retrieved successfully", data)
}
