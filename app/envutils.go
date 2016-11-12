package app

import "os"

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string) bool {
	if env := os.Getenv(key); env == "true" {
		return true
	}
	return false
}
