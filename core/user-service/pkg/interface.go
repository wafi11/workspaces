package pkg

import (
	"context"

	v1 "github.com/wafi11/workspaces/core/user-service/gen/v1"
)

type IRepository interface {
	GetProfile(c context.Context, req *v1.GetProfileRequest) (*v1.GetProfileResponse, error)
	UpdateProfile(c context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error)
}
