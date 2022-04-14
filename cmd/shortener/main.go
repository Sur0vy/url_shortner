package main

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/server"
)

func main() {
	config.Params = *config.SetupConfig(false)
	server.SetupServer().Run(config.Params.ServerAddress)
}
