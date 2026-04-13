package utils

import "fmt"

func GenerateUrl(userID, domain string) string {
	return fmt.Sprintf("%s.%s", userID, domain)
}
