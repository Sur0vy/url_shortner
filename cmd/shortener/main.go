package main

import (
	"database/sql"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/database"
	"github.com/Sur0vy/url_shortner.git/internal/server"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
)

func main() {
	config.Cnf = *config.Setup(config.LoadParams())
	var storage *storage.Storage
	if config.Cnf.DatabaseDSN != "" {
		storage = setupDBStorage()
		defer database.DB.Close()
	} else {
		storage = setupMapStorage()
	}
	server.SetupServer(storage).Run(config.Cnf.ServerAddress)
}

func setupMapStorage() *storage.Storage {
	storage := storage.NewMapStorage()
	return &storage
}

func setupDBStorage() *storage.Storage {
	var err error
	database.DB, err = sql.Open("pgx", config.Cnf.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	//create db connection
	st := storage.NewDBStorage(database.DB)

	//load users from database
	return &st
}
