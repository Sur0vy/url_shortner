package main

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/server"
)

func main() {
	config.Cnf = *config.Setup(config.LoadParams())
	server.SetupServer().Run(config.Cnf.ServerAddress)
}
