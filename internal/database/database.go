package database

import (
	"context"
	"database/sql"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("pgx", config.Cnf.DatabaseDSN)
	if err != nil {
		panic(err)
	}
}

func CreateTables() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS \"url\" (\"id\" SERIAL PRIMARY KEY, \"full\" TEXT UNIQUE,\"short\" TEXT UNIQUE);")
	if err != nil {
		panic(err)
	}

	_, err = DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS \"user\" (\"id\" SERIAL PRIMARY KEY, \"name\" TEXT NOT NULL, \"hash\" TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	_, err = DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS \"url_user\" (\"url\" INTEGER NOT NULL, \"user\" INTEGER NOT NULL, PRIMARY KEY (\"url\", \"user\"));"+
		"ALTER TABLE \"url_user\" DROP CONSTRAINT IF EXISTS \"fk_url_user__url\";"+
		"ALTER TABLE \"url_user\" ADD CONSTRAINT \"fk_url_user__url\" FOREIGN KEY (\"url\") REFERENCES \"url\" (\"id\");"+
		"ALTER TABLE \"url_user\" DROP CONSTRAINT IF EXISTS \"fk_url_user__user\";"+
		"ALTER TABLE \"url_user\" ADD CONSTRAINT \"fk_url_user__user\" FOREIGN KEY (\"user\") REFERENCES \"user\" (\"id\")")
	if err != nil {
		panic(err)
	}
}
