package utils

func GetEnvString(envVars map[string]string, key string) string {
	v, _ := envVars[key]
	return v
}
