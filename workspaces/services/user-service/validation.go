package userservices

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func validateMaxLength(field, value string, max int) error {
	if len(value) > max {
		return fmt.Errorf("%s must not exceed %d characters", field, max)
	}
	return nil
}

func validateUserID(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("email is invalid")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}
