package config

import (
	"sync"
	"time"
)

type EnvConf struct {
	AppHttpPort int
	AppName     string
	AppVersion  string

	LogLevel string

	DbDriver string
	DbDsn    string

	VaultAddr         string
	VaultToken        string
	VaultTransitMount string
	VaultTransitKey   string

	RefreshTokenPepper string

	JwtIssuer         string
	JwksCacheTtl      time.Duration
	ClockSkew         time.Duration
	AccessTokenTtl    time.Duration
	RefreshTokenTtl   time.Duration
	KeyRetireInterval time.Duration

	NatsServerUrl        string
	NatsMaxReconnects    int
	NatsReconnectWaiting time.Duration
}

var (
	instance *EnvConf
	once     sync.Once
)

// Singleton
func GetEnvConf() *EnvConf {
	once.Do(func() {
		instance = load()
	})
	return instance
}
