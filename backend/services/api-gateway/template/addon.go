package template

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/models"
	"github.com/wafi11/workspaces/pkg/response"
)

func (h *Handler) CreateTemplateAddon(c echo.Context) error {
	templateId := c.Param("template_id")
	var req models.CreateAddonRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.svc.CreateTemplateAddon(c.Request().Context(), &req, templateId); err != nil {
		fmt.Printf("Error creating template addon: %v\n", err)
		return response.Error(c, http.StatusInternalServerError, "Failed to create template addon", nil)
	}
	return response.Success(c, http.StatusCreated, "Template addon created successfully", nil)
}

func (h *Handler) GetTemplateAddons(c echo.Context) error {
	templateId := c.Param("template_id")
	addons, err := h.svc.GetTemplateAddons(c.Request().Context(), templateId)
	if err != nil {
		fmt.Printf("Error retrieving template addons: %v\n", err)
		return response.Error(c, http.StatusInternalServerError, "Failed to retrieve template addons", nil)
	}
	return response.Success(c, http.StatusOK, "Successfully retrieved template addons", addons)
}

func (h *Handler) UpdateTemplateAddon(c echo.Context) error {
	id := c.Param("id")
	var req models.CreateAddonRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.svc.UpdateTemplateAddon(c.Request().Context(), id, &req); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update template addon", nil)
	}
	return response.Success(c, http.StatusOK, "Template addon updated successfully", nil)
}

func (h *Handler) DeleteTemplateAddon(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.DeleteTemplateAddon(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete template addon", nil)
	}
	return response.Success(c, http.StatusOK, "Template addon deleted successfully", nil)
}