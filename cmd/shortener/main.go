package main

import (
	"context"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/server"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
)

func main() {
	config.Cnf = *config.Setup(config.LoadParams())
	var dataStorage *storage.Storage
	if config.Cnf.DatabaseDSN != "" {
		dataStorage = setupDBStorage()
	} else {
		dataStorage = setupMapStorage()
	}
	_ = server.SetupServer(dataStorage).Run(config.Cnf.ServerAddress)
}

func setupMapStorage() *storage.Storage {
	st := storage.NewMapStorage()
	return &st
}

func setupDBStorage() *storage.Storage {
	st := storage.NewDBStorage(context.Background())
	return &st
}
