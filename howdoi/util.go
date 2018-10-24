package howdoi

import "os"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getMapDef(m map[string]string, key, fallback string) string {
	if value := m[key]; len(value) > 0 {
		return value
	}
	return fallback
}
