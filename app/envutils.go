package app

import "os"

func envString(key, def string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return def
}

func envBool(key string, def bool) bool {
	if env := os.Getenv(key); env == "true" {
		return true
	}
	return def
}
