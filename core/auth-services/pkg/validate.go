package pkg

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
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
	ErrUnauthorized        = errors.New("unauthorized")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrEmailRequired
	}
	if !emailRegex.MatchString(email) {
		return ErrEmailInvalid
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return ErrPasswordWeak
	}
	return nil
}

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return ErrUsernameRequired
	}
	if len(username) < 3 {
		return ErrUsernameTooShort
	}
	return nil
}
