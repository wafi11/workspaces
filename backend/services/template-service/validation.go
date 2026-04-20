package templateservice

import (
	"errors"
	"fmt"

	"github.com/wafi11/workspaces/pkg/models"
)

func validateListTemplates(req *models.ListTemplatesRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	return nil
}

func validateGetTemplate(req *models.GetTemplateRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if req.TemplateId == "" {
		return errors.New("template_id is required")
	}
	return nil
}
func validateCreateTemplate(req *models.CreateTemplateRequest) error {
	if req == nil {
		return fmt.Errorf("%w: request is required", models.ErrTemplateValidation)
	}
	if req.Name == "" {
		return fmt.Errorf("%w: name is required", models.ErrTemplateValidation)
	}
	if req.Icon == "" {
		return fmt.Errorf("%w: icon is required", models.ErrTemplateValidation)
	}
	for i, v := range req.Variables {
		if v.Key == "" {
			return fmt.Errorf("%w: variable[%d] key is required", models.ErrTemplateValidation, i)
		}
	}
	for i, a := range req.Addons {
		if a.Name == "" {
			return fmt.Errorf("%w: addon[%d] name is required", models.ErrTemplateValidation, i)
		}
		if a.Image == "" {
			return fmt.Errorf("%w: addon[%d] image is required", models.ErrTemplateValidation, i)
		}
	}
	for i, a := range req.Files {
		if a.Filename == "" {
			return fmt.Errorf("%w: file[%d] filename is required", models.ErrTemplateValidation, i)
		}
		if a.SortOrder < 0 {
			return fmt.Errorf("%w: file[%d] sort order is required", models.ErrTemplateValidation, i)
		}
	}
	return nil
}
