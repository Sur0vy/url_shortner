package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
	"strings"
)

var Params EnvConfig

type EnvConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	StoragePath   string `env:"FILE_STORAGE_PATH"`
}

func newEnvConfig() *EnvConfig {
	return &EnvConfig{
		ServerAddress: HostAddr + ":" + HostPort,
		BaseURL:       HostAddr + ":" + HostPort,
		StoragePath:   "",
	}
}

type ConfigOption func(*EnvConfig)

func NotDefaultParams() ConfigOption {
	return func(c *EnvConfig) {
		c.readParams()
		env.Parse(c)
	}
}

func SetupConfig(opts ...ConfigOption) *EnvConfig {
	c := newEnvConfig()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *EnvConfig) readParams() {
	args := os.Args
	fmt.Printf("All arguments: %v\n", args)
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "host to listen on")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "base address")
	flag.StringVar(&c.StoragePath, "f", c.StoragePath, "path to storage file")
	flag.Parse()
}

func (c *EnvConfig) BasePort() string {
	part := strings.Split(c.BaseURL, ":")
	cnt := len(part)
	if cnt > 1 {
		return part[cnt-1]
	} else {
		return HostPort
	}
}
