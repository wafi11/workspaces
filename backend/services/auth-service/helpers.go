package authservices

import (
	"fmt"
)

func GenerateNamespace(userID string) string {
	return fmt.Sprintf("ws-%s", userID)
}
