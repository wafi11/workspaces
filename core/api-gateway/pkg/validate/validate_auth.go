package validate

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
	"google.golang.org/grpc/status"
)

var (
	ErrEmailRequired       = errors.New("email is required")
	ErrEmailInvalid        = errors.New("email format is invalid")
	ErrPasswordRequired    = errors.New("password is required")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrPasswordWeak        = errors.New("password must contain uppercase, lowercase, and a number")
	ErrUsernameRequired    = errors.New("username is required")
	ErrUsernameTooShort    = errors.New("username must be at least 3 characters")
	ErrEmailAlreadyExist   = errors.New("email already exists")
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)

func HandleAuthError(c echo.Context, err error) error {
	if st, ok := status.FromError(err); ok {
		msg := st.Message()
		switch msg {
		case ErrEmailAlreadyExist.Error():
			return pkg.Error(c, 400, "Email already exists", nil)
		case ErrPasswordWeak.Error(), ErrPasswordTooShort.Error(), ErrPasswordRequired.Error():
			return pkg.Error(c, 400, msg, nil)
		case ErrEmailRequired.Error(), ErrEmailInvalid.Error():
			return pkg.Error(c, 400, msg, nil)
		case ErrInvalidCredentials.Error():
			return pkg.Error(c, 401, "Invalid credentials", nil)
		default:
			log.Printf("error %s", err.Error())
			return pkg.Error(c, http.StatusInternalServerError, "Internal Server Error", nil)
		}
	}
	return pkg.Error(c, http.StatusInternalServerError, "Internal Server Error", nil)
}
