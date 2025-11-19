package env

import (
	"fmt"
	"log/slog"
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

func (e Env) IsLocal() bool {
	return e == Local
}

func (e Env) IsDev() bool {
	return e == Dev
}

func (e Env) IsStg() bool {
	return e == Stg
}

func (e Env) IsProd() bool {
	return e == Prod
}

func GetString(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Warn(fmt.Sprintf("key %s not found, using fallback %s", key, fallback))
		return fallback
	}

	return value
}

func GetInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		slog.Warn(fmt.Sprintf("key %s not found, using fallback %d", key, fallback))
		return fallback
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		slog.Warn(fmt.Sprintf("key %s not found, using fallback %d", key, fallback))
		return fallback
	}

	return i
}

func GetBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		slog.Warn(fmt.Sprintf("key %s not found, using fallback %t", key, fallback))
		return fallback
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		slog.Warn(fmt.Sprintf("key %s not found, using fallback %t", key, fallback))
		return fallback
	}

	return b
}
