package config

import "time"

func GetDebug(getenv string) bool {
	if getenv == "" {
		return false
	}

	return getenv == "true"
}

func GetEnv(getenv string, address string) string {
	if getenv == "" {
		return address
	}

	return getenv
}

func ParseDuration(getenv string, timeout time.Duration) time.Duration {
	if getenv == "" {
		return timeout
	}

	d, err := time.ParseDuration(getenv)
	if err != nil {
		return timeout
	}

	return d
}
