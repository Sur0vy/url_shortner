package config

import (
	"github.com/caarlos0/env/v6"
)

var Params EnvConfig

type EnvConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func newEnvConfig() *EnvConfig {
	return &EnvConfig{
		ServerAddress: "",
		BaseURL:       "",
	}
}

func SetupConfig() *EnvConfig {
	c := newEnvConfig()
	err := env.Parse(&c)
	if err != nil {
		c.ServerAddress = ServerAddress
		c.BaseURL = HTTPPref
	}
	return c
}
