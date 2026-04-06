package workspace

import (
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
	username := c.Get("username").(string)
	userId := c.Get("user_id").(string)
	var req workspaceservice.CreateWorkspaceRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid body request", nil)
	}

	req.UserId = userId

	data, err := h.svc.CreateWorkspace(c.Request().Context(), &req, username)
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

func (h *Handler) DetailsWorkspaces(c echo.Context) error {
   workspaceId := c.Param("id")

    data, err := h.svc.GetWorkspace(c.Request().Context(), &workspaceservice.GetWorkspaceRequest{
		WorkspaceId: workspaceId,
		UserId: "",
	})
    if err != nil {
        return response.Error(c, http.StatusInternalServerError, "failed to list workspaces", err)
    }

    return response.Success(c, http.StatusOK, "Successfullt Get Workspaces", data.Workspace)
}
