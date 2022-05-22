package users

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Sur0vy/url_shortner.git/internal/storage/helpers"
)

type DBUserStorage struct {
	database *sql.DB
}

func NewDBUserStorage(db *sql.DB) UserStorage {
	return &DBUserStorage{
		database: db,
	}
}

func (u *DBUserStorage) Add(ctx context.Context) (string, string) {
	user := User + strconv.Itoa(u.GetCount(ctx)+1)
	hash := GenerateHash(user)

	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		helpers.TUser, helpers.FUserName, helpers.FUserHash)
	_, err := u.database.ExecContext(ctxIn, sqlStr, user, hash)
	if err != nil {
		return "", ""
	}

	return user, hash
}

func (u *DBUserStorage) GetUser(ctx context.Context, hash string) string {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var name string

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		helpers.FUserName, helpers.TUser, helpers.FUserHash)
	row := u.database.QueryRowContext(ctxIn, sqlStr, hash)
	err := row.Scan(&name)

	if err != nil {
		return ""
	}

	return name
}

func (u *DBUserStorage) LoadFromFile() error {
	return nil
}

func (u *DBUserStorage) GetCount(ctx context.Context) int {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var ret int
	sqlStr := fmt.Sprintf("SELECT count(%s) FROM %s",
		helpers.FUserName, helpers.TUser)
	row := u.database.QueryRowContext(ctxIn, sqlStr)
	err := row.Scan(&ret)

	if err != nil {
		return 0
	}
	return ret
}
