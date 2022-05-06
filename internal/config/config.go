package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
	"strings"
)

var Cnf Config

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	StoragePath   string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
	CurrentUser   string
}

func newConfig() *Config {
	return &Config{
		ServerAddress: HostAddr + ":" + HostPort,
		BaseURL:       HostAddr + ":" + HostPort,
		StoragePath:   "",
	}
}

type option func(*Config)

func LoadParams() option {
	return func(c *Config) {
		c.flagParams()
		env.Parse(c)
	}
}

func Setup(opts ...option) *Config {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Config) flagParams() {
	args := os.Args
	fmt.Printf("All arguments: %v\n", args)
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "host to listen on")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "base address")
	flag.StringVar(&c.StoragePath, "f", c.StoragePath, "path to storage file")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "address to database connection")
	flag.Parse()
}

func (c *Config) BasePort() string {
	part := strings.Split(c.BaseURL, ":")
	cnt := len(part)
	if cnt > 1 {
		return part[cnt-1]
	} else {
		return HostPort
	}
}
