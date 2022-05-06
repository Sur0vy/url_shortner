package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB
