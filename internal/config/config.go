package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	loadfile "github.com/safakaveh/card-platform/internal/common/load-file"
)

func load() *EnvConf {
	loadfile.LoadEnvFile()
	cfg := &EnvConf{}

	cfg.AppHttpPort = getEnvAsInt("APP_HTTP_PORT", 8080)

	cfg.AppVersion = getEnv("APP_VERSION", "2.0.0")

	cfg.AppName = getEnv("APP_NAME", "card-platform")

	cfg.LogLevel = getEnv("LOG_LEVEL", "info")

	cfg.DbDriver = getEnv("DB_DRIVER", "sqlite")

	cfg.DbDatasourceName = getEnv("DB_DATASOURCE_NAME", "./data.db")

	return cfg
}

func getEnv(key, defaultVal string) string {
	key = strings.TrimSpace(key)
	if val := os.Getenv(key); val != "" {
		return val
	}
	return strings.TrimSpace(defaultVal)
}

func mustGetEnv(key string) string {
	key = strings.TrimSpace(key)
	val := os.Getenv(key)
	if val == "" {
		panic("missing required env: " + key)
	}
	return strings.TrimSpace(val)
}

func getEnvAsInt(key string, defaultVal int) int {
	key = strings.TrimSpace(key)
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	key = strings.TrimSpace(key)
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
