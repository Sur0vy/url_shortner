package config

import (
	"github.com/caarlos0/env/v6"
)

var Params EnvConfig

type EnvConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	StoragePath   string `env:"FILE_STORAGE_PATH"`
}

func newEnvConfig() *EnvConfig {
	return &EnvConfig{
		ServerAddress: ServerAddress,
		BaseURL:       HTTPPref,
		StoragePath:   "",
	}
}

func SetupConfig() *EnvConfig {
	c := newEnvConfig()
	env.Parse(c)
	return c
}
