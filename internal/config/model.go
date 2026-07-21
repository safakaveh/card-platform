package config

import (
	"sync"
)

type EnvConf struct {
	AppHttpPort string
	AppName     string
	AppVersion  string

	LogLevel string

	DbDriver         string
	dbDatasourceName string
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
