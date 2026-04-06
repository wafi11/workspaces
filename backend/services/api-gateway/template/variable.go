package template

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/response"
	templateservice "github.com/wafi11/workspaces/services/template-service"
)

func (h *Handler) CreateTemplateVariable(c echo.Context) error {
	templateId := c.Param("template_id")
	var req templateservice.CreateVariableRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.svc.CreateTemplateVariable(c.Request().Context(), &req, templateId); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create template variable", nil)
	}
	return response.Success(c, http.StatusCreated, "Template variable created successfully", nil)
}

func (h *Handler) GetTemplateVariables(c echo.Context) error {
	templateId := c.Param("template_id")
	variables, err := h.svc.GetTemplateVariables(c.Request().Context(), templateId)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to retrieve template variables", nil)
	}
	return response.Success(c, http.StatusOK, "Successfully retrieved template variables", variables)
}


func (h *Handler) UpdateTemplateVariable(c echo.Context) error {
	id := c.Param("id")
	var req templateservice.CreateVariableRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.svc.UpdateTemplateVariable(c.Request().Context(), id, &req); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update template variable", nil)
	}
	return response.Success(c, http.StatusOK, "Template variable updated successfully", nil)
}

func (h *Handler) DeleteTemplateVariable(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.DeleteTemplateVariable(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete template variable", nil)
	}
	return response.Success(c, http.StatusOK, "Template variable deleted successfully", nil)
}