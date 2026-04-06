package templateservice

import (
	"errors"
	"fmt"
)

func validateListTemplates(req *ListTemplatesRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	return nil
}

func validateGetTemplate(req *GetTemplateRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if req.TemplateId == "" {
		return errors.New("template_id is required")
	}
	return nil
}
func validateCreateTemplate(req *CreateTemplateRequest) error {
	if req == nil {
		return fmt.Errorf("%w: request is required", ErrValidation)
	}
	if req.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if req.Image == "" {
		return fmt.Errorf("%w: image is required", ErrValidation)
	}
	for i, v := range req.Variables {
		if v.Key == "" {
			return fmt.Errorf("%w: variable[%d] key is required", ErrValidation, i)
		}
	}
	for i, a := range req.Addons {
		if a.Name == "" {
			return fmt.Errorf("%w: addon[%d] name is required", ErrValidation, i)
		}
		if a.Image == "" {
			return fmt.Errorf("%w: addon[%d] image is required", ErrValidation, i)
		}
	}
	for i, a := range req.Files {
		if a.Filename == "" {
			return fmt.Errorf("%w: file[%d] filename is required", ErrValidation, i)
		}
		if a.SortOrder < 0 {
			return fmt.Errorf("%w: file[%d] sort order is required", ErrValidation, i)
		}
	}
	return nil
}
