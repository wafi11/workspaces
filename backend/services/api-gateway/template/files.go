package template

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	templateservice "github.com/wafi11/workspaces/services/template-service"
)

func (h *Handler) CreateTemplateFiles(c echo.Context) error {
	templateId := c.Param("template_id")
	var req templateservice.CreateTemplateFilesRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.svc.CreateTemplateFiles(c.Request().Context(), &req, templateId); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create template files", nil)
	}
	return response.Success(c, http.StatusCreated, "Template files created successfully", nil)
}


func (h *Handler) GetTemplateFiles(c echo.Context) error {
	templateId := c.Param("template_id")
	files, err := h.svc.GetTemplateFiles(c.Request().Context(), templateId)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to retrieve template files", nil)
	}
	return response.Success(c, http.StatusOK, "Successfully retrieved template files", files)
}

func (h *Handler) UpdateTemplateFiles(c echo.Context) error {
	id := c.Param("id")
	var req templateservice.CreateTemplateFilesRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.svc.UpdateTemplateFiles(c.Request().Context(), id, &req); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update template files", nil)
	}

	return response.Success(c, http.StatusOK, "Template files updated successfully", nil)
}

func (h *Handler) DeleteTemplateFiles(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.DeleteTemplateFiles(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete template files", nil)
	}
	return response.Success(c, http.StatusOK, "Template files deleted successfully", nil)
}