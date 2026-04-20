package template

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/pkg/models"
	"github.com/wafi11/workspaces/pkg/response"
	templateservice "github.com/wafi11/workspaces/services/template-service"
)

type Handler struct {
	svc templateservice.TemplateService
}

func NewHandler(svc templateservice.TemplateService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateTemplate(ctx echo.Context) error {
	var req models.CreateTemplateRequest

	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
	}

	data, err := h.svc.CreateTemplate(ctx.Request().Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTemplateValidation):
			return response.Error(ctx, http.StatusBadRequest, err.Error(), err)
		case errors.Is(err, models.ErrTemplateNotFound):
			return response.Error(ctx, http.StatusNotFound, "template not found", err)
		default:
			return response.Error(ctx, http.StatusInternalServerError, "failed to create template", err)
		}
	}

	return response.Success(ctx, http.StatusCreated, "template created successfully", data)
}

func (h *Handler) FindTemplateDetailsForm(c echo.Context) error {
	template_id := c.Param("template_id")

	if template_id == "" {
		return response.Error(c, http.StatusNotFound, "templates not found", nil)
	}

	data, err := h.svc.GetDetailsInfo(c.Request().Context(), template_id)

	if err != nil {
		return response.Error(c, http.StatusNotFound, "templates not found", nil)
	}

	return response.Success(c, http.StatusOK, "Successfully retreived templates details", data)
}

func (h *Handler) UpdateTemplate(ctx echo.Context) error {
	id := ctx.Param("id")
	var req models.UpdateTemplateRequest

	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
	}

	err := h.svc.UpdateTemplate(ctx.Request().Context(), id, &req)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTemplateValidation):
			return response.Error(ctx, http.StatusBadRequest, err.Error(), err)
		case errors.Is(err, models.ErrTemplateNotFound):
			return response.Error(ctx, http.StatusNotFound, "template not found", err)
		default:
			return response.Error(ctx, http.StatusInternalServerError, "failed to create template", err)
		}
	}

	return response.Success(ctx, http.StatusCreated, "Template update successfully", nil)
}

func (h *Handler) DeleteTemplate(ctx echo.Context) error {
	id := ctx.Param("id")

	err := h.svc.DeleteTemplate(ctx.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTemplateValidation):
			return response.Error(ctx, http.StatusBadRequest, err.Error(), err)
		case errors.Is(err, models.ErrTemplateNotFound):
			return response.Error(ctx, http.StatusNotFound, "template not found", err)
		default:
			return response.Error(ctx, http.StatusInternalServerError, "failed to create template", err)
		}
	}

	return response.Success(ctx, http.StatusCreated, "Template delete successfully", nil)
}

func (h *Handler) GetTemplateDetails(ctx echo.Context) error {
	template_id := ctx.Param("template_id")

	data, err := h.svc.GetTemplate(ctx.Request().Context(), &models.GetTemplateRequest{
		TemplateId: template_id,
	})

	if err != nil {
		return response.Error(ctx, http.StatusNotFound, "template not found", nil)
	}
	log.Printf("[template] id=%s name=%s icon=%s url=%s",
		data.Template.Id,
		data.Template.Name,
		data.Template.Icon,
		data.Template.TemplateUrl,
	)

	return response.Success(ctx, http.StatusOK, "Successfully Get Template", data.Template)

}

func (h *Handler) GetListTemplates(ctx echo.Context) error {
	req := &models.ListTemplatesRequest{
		Category: ctx.QueryParam("category"),
	}

	data, err := h.svc.ListTemplates(ctx.Request().Context(), req)
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "failed to fetch templates", err)
	}
	

	return response.Success(ctx, http.StatusOK, "successfully get templates", data.Templates)
}


func (h *Handler) FinTemplateWorkspaceForm(ctx echo.Context) error {

	data, err := h.svc.FindTemplateWorkspaceForm(ctx.Request().Context())
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, err.Error(),nil)
	}

	return response.Success(ctx, http.StatusOK, "successfully get templates workspace form", data)
}
