package workspaceservice

import (
	"errors"
	"fmt"
	"regexp"
)

func validateCreateWorkspace(req *CreateWorkspaceRequest) error {
	if req == nil {
		return fmt.Errorf("%w: request is required", ErrValidation)
	}

	switch {
	case req.UserId == "":
		return fmt.Errorf("%w: user_id is required", ErrValidation)
	case req.TemplateId == "":
		return fmt.Errorf("%w: template_id is required", ErrValidation)
	}

	if err := validateK8sName(req.Name); err != nil {
		return err
	}

	if _, err := validationTypeSchedulling(req.TypeTimeDuration); err != nil {
		return fmt.Errorf("%w: %s", ErrValidation, err)
	}

	return nil
}

var k8sNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

func validateK8sName(name string) error {
	if len(name) < 2 {
		return fmt.Errorf("%w: name must be at least 2 characters", ErrValidation)
	}
	if len(name) > 63 {
		return fmt.Errorf("%w: name must be 63 characters or less", ErrValidation)
	}
	if !k8sNameRegex.MatchString(name) {
		return fmt.Errorf("%w: name must be lowercase alphanumeric and hyphens only, and must start and end with alphanumeric", ErrValidation)
	}
	return nil
}
func validateListWorkspaces(req *ListWorkspacesRequest) error {
	switch {
	case req == nil:
		return errors.New("request is required")
	case req.UserId == "":
		return fmt.Errorf("%w: user_id is required", ErrValidation)
	}
	return nil
}

func validateGetWorkspace(req *GetWorkspaceRequest, isAdmin bool) error {
	switch {
	case req == nil:
		return errors.New("request is required")
	case req.WorkspaceId == "":
		return fmt.Errorf("%w: workspace_id is required", ErrValidation)
	case !isAdmin && req.UserId == "":
		return fmt.Errorf("%w: user_id is required", ErrValidation)
	}
	return nil
}

func validateDeleteWorkspace(req *DeleteWorkspaceRequest) error {
	switch {
	case req == nil:
		return errors.New("request is required")
	case req.WorkspaceId == "":
		return fmt.Errorf("%w: workspace_id is required", ErrValidation)
	case req.UserId == "":
		return fmt.Errorf("%w: user_id is required", ErrValidation)
	}
	return nil
}
