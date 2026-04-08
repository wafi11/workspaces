package workspaceservice

import (
	"errors"
	"fmt"
)


func validateCreateWorkspace(req *CreateWorkspaceRequest) error {
	switch {
	case req == nil:
		return fmt.Errorf("%w: request is required", ErrValidation)
	case req.UserId == "":
		return fmt.Errorf("%w: user_id is required", ErrValidation)
	case req.TemplateId == "":
		return fmt.Errorf("%w: template_id is required", ErrValidation)
	case req.Name == "":
		return fmt.Errorf("%w: name is required", ErrValidation)
	case len(req.Name) > 50:
		return fmt.Errorf("%w: name must be 50 characters or less", ErrValidation)
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