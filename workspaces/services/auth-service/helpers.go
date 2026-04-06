package authservices

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateNamespace(email string) string {
	local := strings.Split(email, "@")[0]
	local = strings.ToLower(local)
	local = strings.ReplaceAll(local, ".", "-")
	local = strings.ReplaceAll(local, "_", "-")
	suffix := uuid.New().String()[:8]
	return fmt.Sprintf("user-%s-%s", local, suffix)
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashed), nil
}

func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
