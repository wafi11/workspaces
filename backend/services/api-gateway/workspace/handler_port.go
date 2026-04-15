package workspace

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
)

var allowedPorts = map[int]bool{
	3000: true, 3001: true, 3002: true, 3003: true, 3004: true,
	3005: true, 3006: true, 3007: true, 3008: true, 3009: true, 3010: true,
}

func (h *Handler) CreateWorkspacePort(e echo.Context) error {
	workspaceID := e.Param("workspace_id")
	userID := e.Get("user_id").(string)

	if workspaceID == "" {
		return response.Error(e, http.StatusBadRequest, "workspace id not found", nil)
	}

	var req struct {
		Port int `json:"port"`
	}
	if err := e.Bind(&req); err != nil {
		return response.Error(e, http.StatusBadRequest, "invalid request body", err)
	}

	if !allowedPorts[req.Port] {
		return response.Error(e, http.StatusBadRequest, "port must be between 3000-3010", nil)
	}

	// validasi workspace milik user
	if err := h.svc.ValidateWorkspaceOwner(e.Request().Context(), workspaceID, userID); err != nil {
		return response.Error(e, http.StatusForbidden, "workspace not found or unauthorized", nil)
	}

	_, err := h.svc.CreateWorkspacePort(e.Request().Context(), workspaceID, req.Port, userID)
	if err != nil {
		log.Printf("port erro : %s",err.Error())
		return response.Error(e, http.StatusBadRequest, err.Error(), err)
	}

	return response.Success(e, http.StatusCreated, "Successfully Create Port", nil)
}

func (h *Handler) DeleteWorkspacePort(e echo.Context) error {
	workspaceID := e.Param("workspace_id")
	userID := e.Get("user_id").(string)

	if workspaceID == "" {
		return response.Error(e, http.StatusBadRequest, "workspace id not found", nil)
	}

	var req struct {
		Port int `json:"port"`
	}
	if err := e.Bind(&req); err != nil {
		return response.Error(e, http.StatusBadRequest, "invalid request body", err)
	}

	if !allowedPorts[req.Port] {
		return response.Error(e, http.StatusBadRequest, "port must be between 3000-3010", nil)
	}

	if err := h.svc.ValidateWorkspaceOwner(e.Request().Context(), workspaceID, userID); err != nil {
		return response.Error(e, http.StatusForbidden, "workspace not found or unauthorized", nil)
	}

	if err := h.svc.DeleteWorkspacePort(e.Request().Context(), workspaceID, req.Port); err != nil {
		return response.Error(e, http.StatusBadRequest, err.Error(), err)
	}

	return response.Success(e, http.StatusOK, "Successfully Delete Port", nil)
}

func (h *Handler) ListWorkspacePorts(e echo.Context) error {
	workspaceID := e.Param("workspace_id")
	userID := e.Get("user_id").(string)

	if workspaceID == "" {
		return response.Error(e, http.StatusBadRequest, "workspace id not found", nil)
	}

	if err := h.svc.ValidateWorkspaceOwner(e.Request().Context(), workspaceID, userID); err != nil {
		return response.Error(e, http.StatusForbidden, "workspace not found or unauthorized", nil)
	}

	ports, err := h.svc.ListWorkspacePorts(e.Request().Context(), workspaceID)
	if err != nil {
		return response.Error(e, http.StatusInternalServerError, err.Error(), err)
	}

	return response.Success(e, http.StatusOK, "Successfully Get Ports", ports)
}