package utils

func GetEnvString(envVars map[string]any, key string) string {
	v, _ := envVars[key].(string)
	return v
}
