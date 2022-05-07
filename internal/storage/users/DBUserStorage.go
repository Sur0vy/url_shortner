package users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/database"
	"strconv"
	"time"
)

type DBUserStorage struct {
	database *sql.DB
}

func NewDBUserStorage(db *sql.DB) UserStorage {
	return &DBUserStorage{
		database: db,
	}
}

func (u *DBUserStorage) Add() (string, string) {
	user := User + strconv.Itoa(u.GetCount()+1)
	hash := GenerateHash(user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		database.TUser, database.FUserName, database.FUserHash)
	_, err := u.database.ExecContext(ctx, sql, user, hash)
	if err != nil {
		return "", ""
	}

	return user, hash
}

func (u *DBUserStorage) GetUser(hash string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var name string

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		database.FUserName, database.TUser, database.FUserHash)
	row := u.database.QueryRowContext(ctx, sql, hash)
	err := row.Scan(&name)

	if err != nil {
		return ""
	}

	return name
}

func (u *DBUserStorage) LoadFromFile() error {
	return nil
}

func (u *DBUserStorage) GetCount() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ret int
	sql := fmt.Sprintf("SELECT count(%s) FROM %s",
		database.FUserName, database.TUser)
	row := u.database.QueryRowContext(ctx, sql)
	err := row.Scan(&ret)

	if err != nil {
		return 0
	}
	return ret
}
