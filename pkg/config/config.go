package config

import (
	"strings"
	"time"
)

func GetEnvBoolean(getenv string) bool {
	if getenv == "" {
		return false
	}

	return getenv == "true"
}

func GetEnv(getenv string, def string) string {
	if getenv == "" {
		return def
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

func GetBrokers(getenv string, def []string) []string {
	if getenv == "" {
		return def
	}

	return strings.Split(getenv, ",")
}
