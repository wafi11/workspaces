package workspace

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	workspaceservice "github.com/wafi11/workspaces/services/workspaces-service"
)

type Handler struct {
	svc workspaceservice.WorkspaceService
}

func NewHandler(svc workspaceservice.WorkspaceService) *Handler {
	return &Handler{
		svc: svc,
	}
}



func (h *Handler) Create(c echo.Context) error {
	userId := c.Get("user_id").(string)
	var req workspaceservice.CreateWorkspaceRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid body request", nil)
	}

	req.UserId = userId

	data, err := h.svc.CreateWorkspace(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to create workspace", err)
	}

	return response.Success(c, http.StatusCreated, "Sucessfully create workspaces", data)
}

func (h *Handler) ListWorkspaces(c echo.Context) error {
	limitInt, _ := strconv.Atoi(c.QueryParam("limit"))
	offsetInt, _ := strconv.Atoi(c.QueryParam("offset"))
	status := c.QueryParam("status")

	if limitInt <= 0 || limitInt > 100 {
		limitInt = 20
	}
	if offsetInt < 0 {
		offsetInt = 0
	}

	data, err := h.svc.ListWorkspaces(c.Request().Context(), (limitInt), offsetInt, status)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to list workspaces", err)
	}

	return response.Success(c, http.StatusOK, "successfully list workspaces", data.Workspaces)
}

func (h *Handler) ListWorkspaceForm(c echo.Context) error {
	userId := c.Get("user_id").(string)
	if userId == "" {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	data, err := h.svc.ListWorkspaceForm(c.Request().Context(), userId)

	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	return response.Success(c, http.StatusOK, "Successfully Retreived workspace form", data)
}

func (h *Handler) ListWorkspaceByUserId(c echo.Context) error {
	userId := c.Get("user_id").(string)

	if userId == "" {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}
	data, err := h.svc.ListWorkspacesByUserId(c.Request().Context(), &workspaceservice.ListWorkspacesRequest{
		UserId: userId,
	})

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "inetrnnla server error", nil)
	}

	return response.Success(c, http.StatusOK, "Successfully retreived data workspaces", data.Workspaces)
}

func (h *Handler) DetailsWorkspaces(c echo.Context) error {
	workspaceId := c.Param("id")

	data, err := h.svc.GetWorkspace(c.Request().Context(), &workspaceservice.GetWorkspaceRequest{
		WorkspaceId: workspaceId,
		UserId:      "",
	})
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to list workspaces", err)
	}

	return response.Success(c, http.StatusOK, "Successfullt Get Workspaces", data.Workspace)
}

func (h *Handler) StartWorkspaces(c echo.Context) error {
	workspaceId := c.Param("workspace_id")
	userId := c.Get("user_id").(string)

	if workspaceId == "" {
		return response.Error(c, http.StatusBadRequest, "workspace not found", nil)
	}

	err := h.svc.UpdateWorkspaceStatus(c.Request().Context(), workspaceId,userId, "running")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, "Successfully Start Workspace", nil)
}

func (h *Handler) PausedWorkspaces(c echo.Context) error {
	workspaceId := c.Param("workspace_id")
	userId := c.Get("user_id").(string)

	if workspaceId == "" {
		return response.Error(c, http.StatusBadRequest, "workspace not found", nil)
	}

	err := h.svc.UpdateWorkspaceStatus(c.Request().Context(), workspaceId,userId, "paused")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to paused workspaces", nil)
	}

	return response.Success(c, http.StatusOK, "Successfully Paused Workspace", nil)
}


func (h *Handler) ResumedWorkspaces(c echo.Context) error {
	workspaceId := c.Param("workspace_id")
	userId := c.Get("user_id").(string)

	if workspaceId == "" {
		return response.Error(c, http.StatusBadRequest, "workspace not found", nil)
	}

	err := h.svc.UpdateWorkspaceStatus(c.Request().Context(), workspaceId,userId, "resumed")
	if err != nil {
		log.Printf("error : %s",err.Error())
		return response.Error(c, http.StatusBadRequest, "failed to resumed workspaces", err)
	}

	return response.Success(c, http.StatusOK, "Successfully Resumed Workspace", nil)
}

func (h *Handler) StopWorkspaces(c echo.Context) error {
	workspaceId := c.Param("workspace_id")
	userId := c.Get("user_id").(string)

	if workspaceId == "" {
		return response.Error(c, http.StatusBadRequest, "workspace not found", nil)
	}

	err := h.svc.UpdateWorkspaceStatus(c.Request().Context(), workspaceId,userId, "stopped")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to stop workspaces", nil)
	}

	return response.Success(c, http.StatusOK, "Successfully Stop Workspace", nil)
}
