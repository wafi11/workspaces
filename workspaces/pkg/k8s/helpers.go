package k8s

import "fmt"

func GenerateNamespace(userId, name string) string {
	return fmt.Sprintf("ws-%s-%s", userId[:8], name)
}
