package authservices

import (
	"fmt"
)

func generateNamespace(userID string) string {
	return fmt.Sprintf("ws-%s", userID)
}
