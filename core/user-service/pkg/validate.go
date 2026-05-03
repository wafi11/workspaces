package pkg

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const (
	BucketProfile string = "profiles"
)
