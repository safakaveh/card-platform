package config

import (
	"sync"
)

type EnvConf struct {
	AppHttpPort int
	AppName     string
	AppVersion  string

	LogLevel string

	DbDriver         string
	DbDatasourceName string
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
