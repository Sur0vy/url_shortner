package users

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

type DBUserStorage struct {
	counter  int
	database *sql.DB
}

func NewDBUserStorage(db *sql.DB) UserStorage {
	return &DBUserStorage{
		counter:  0,
		database: db,
	}
}

func (u *DBUserStorage) Add() (string, string) {
	u.counter++
	user := User + strconv.Itoa(u.counter)
	hash := GenerateHash(user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := "INSERT INTO \"user\" (name, hash) VALUES ($1, $2)"
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

	sql := "SELECT name FROM \"user\" WHERE hash = $1"
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
	sql := "SELECT count(name) FROM \"user\""
	row := u.database.QueryRowContext(ctx, sql)
	err := row.Scan(&ret)

	if err != nil {
		return 0
	}
	return ret
}
