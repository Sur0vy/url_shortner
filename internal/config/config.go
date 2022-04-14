package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
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

func SetupConfig(defaultParams bool) *EnvConfig {
	c := newEnvConfig()
	if !defaultParams {
		env.Parse(c)
		c.readParams()
	}
	return c
}

func (c *EnvConfig) readParams() {
	args := os.Args
	fmt.Printf("All arguments: %v\n", args)
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "host to listen on")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "base address")
	flag.StringVar(&c.StoragePath, "f", c.StoragePath, "path to storage file")
	//flag.StringVar(&c.ServerAddress, "a", os.Getenv("SERVER_ADDRESS"), "host to listen on")
	//flag.StringVar(&c.BaseURL, "b", os.Getenv("BASE_URL"), "base address")
	//flag.StringVar(&c.StoragePath, "f", os.Getenv("FILE_STORAGE_PATH"), "path to storage file")
	flag.Parse()
}
