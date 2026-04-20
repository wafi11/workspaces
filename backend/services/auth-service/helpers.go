package authservices

import (
	"fmt"
)

// keys.go atau di atas file yang sama
func refreshTokenKey(sessionID string) string {
	return fmt.Sprintf("refresh_token:%s", sessionID)
}

func sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}
func GenerateNamespace(userID string) string {
	return fmt.Sprintf("ws-%s", userID)
}
