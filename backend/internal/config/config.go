package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	loadfile "github.com/sunecity/smart-building-platform/auth/internal/common/load-file"
)

func load() *EnvConf {
	loadfile.LoadEnvFile()
	cfg := &EnvConf{}

	cfg.AppVersion = getEnv("APP_VERSION", "3.0.0")

	cfg.AppName = getEnv("APP_NAME", "auth")

	cfg.LogLevel = getEnv("LOG_LEVEL", "info")

	cfg.DbDriver = getEnv("DB_DRIVER", "postgres")

	cfg.DbDsn = mustGetEnv("DB_DSN")

	cfg.VaultAddr = getEnv("VAULT_ADDR", "http://127.0.0.1:8200")

	cfg.VaultToken = mustGetEnv("VAULT_TOKEN")

	cfg.VaultTransitMount = getEnv("VAULT_TRANSIT_MOUNT", "transit")

	cfg.VaultTransitKey = getEnv("VAULT_TRANSIT_KEY", "jwt-ed25519")

	cfg.RefreshTokenPepper = mustGetEnv("REFRESH_TOKEN_PEPPER")

	cfg.JwtIssuer = getEnv("JWT_ISSUER", "sun-auth-v3.0.0")

	jwtCacheTtl, _ := time.ParseDuration("5m")
	cfg.JwksCacheTtl = getEnvAsDuration("JWKS_CACHE_TTL", jwtCacheTtl) // 5 minutes

	accessTokenTtl, _ := time.ParseDuration("7m")
	cfg.AccessTokenTtl = getEnvAsDuration("ACCESS_TOKEN_TTL", accessTokenTtl) // 7 minutes

	refreshTokenTtl, _ := time.ParseDuration("7m")
	cfg.RefreshTokenTtl = getEnvAsDuration("REFRESH_TOKEN_TTL", refreshTokenTtl) // 30 days

	clockSkew, _ := time.ParseDuration("60s")
	cfg.ClockSkew = getEnvAsDuration("CLOCK_SKEW", clockSkew) // 60 seconds

	keyRetireInterval, _ := time.ParseDuration("1m")
	cfg.KeyRetireInterval = getEnvAsDuration("KEY_RETIRE_INTERVAL", keyRetireInterval)

	cfg.NatsServerUrl = getEnv("NATS1_SERVER_URL", "nats://127.0.0.1:4222")

	cfg.NatsMaxReconnects = getEnvAsInt("NATS1_MAX_RECONNECTS", 5)

	natsReconnectWaiting, _ := time.ParseDuration("60s")

	cfg.NatsReconnectWaiting = getEnvAsDuration("NATS1_RECONNECT_WAITING", natsReconnectWaiting)
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
