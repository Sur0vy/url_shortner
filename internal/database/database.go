package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

var DB *sql.DB

const (
	TUser     string = "u_user"
	TURL      string = "u_url"
	TUserURL  string = "u_user_URL"
	FShort    string = "URL_short"
	FFull     string = "URL_full"
	FUserName string = "User_name"
	FUserHash string = "User_hash"
)

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

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT UNIQUE, %s TEXT UNIQUE PRIMARY KEY)",
		TURL, FFull, FShort)
	_, err := DB.ExecContext(ctx, sql)
	if err != nil {
		panic(err)
	}

	sql = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT NOT NULL, %s TEXT NOT NULL PRIMARY KEY)",
		TUser, FUserName, FUserHash)
	_, err = DB.ExecContext(ctx, sql)
	if err != nil {
		panic(err)
	}

	sql = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT NOT NULL, %s TEXT NOT NULL, PRIMARY KEY (%s, %s))",
		TUserURL, FShort, FUserHash, FShort, FUserHash)
	_, err = DB.ExecContext(ctx, sql)
	if err != nil {
		panic(err)
	}

	sql = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_u_us_URLID;"+
		"ALTER TABLE %s ADD CONSTRAINT fk_u_us_URLID FOREIGN KEY (%s) REFERENCES %s (%s);",
		TUserURL, TUserURL, FShort, TURL, FShort)
	_, err = DB.ExecContext(ctx, sql)
	if err != nil {
		panic(err)
	}

	sql = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_u_us_UserID;"+
		"ALTER TABLE %s ADD CONSTRAINT fk_u_us_UserID FOREIGN KEY (%s) REFERENCES %s (%s);",
		TUserURL, TUserURL, FUserHash, TUser, FUserHash)
	_, err = DB.ExecContext(ctx, sql)
	if err != nil {
		panic(err)
	}
}
