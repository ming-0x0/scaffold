package env

import (
	"os"
	"strconv"
)

type Env string

const (
	Local Env = "local"
	Dev   Env = "dev"
	Stg   Env = "stg"
	Prod  Env = "prod"
)

func GetEnv() Env {
	switch os.Getenv("ENV") {
	case "local":
		return Local
	case "dev":
		return Dev
	case "stg":
		return Stg
	case "prod":
		return Prod
	default:
		return Local
	}
}

func GetString(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func GetInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return i
}

func GetBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return b
}
