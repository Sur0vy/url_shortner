package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/storage/users"
	_ "github.com/jackc/pgx/v4/stdlib"
	"strconv"
	"time"
)

type DBStorage struct {
	userStorage users.UserStorage
	counter     int
	database    *sql.DB
}

func NewDBStorage(db *sql.DB) Storage {
	s := &DBStorage{
		database:    db,
		counter:     0,
		userStorage: users.NewDBUserStorage(db),
	}
	return s
}

func (s *DBStorage) InsertURL(fullURL string) string {
	s.counter++
	short := strconv.Itoa(s.counter)

	fmt.Printf("\tAdd new URL to storage full = %s, short = %s\n", fullURL, short)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := "INSERT INTO \"url\" (\"full\", \"short\") VALUES ($1, $2)"
	_, err := s.database.ExecContext(ctx, sql, fullURL, short)

	if err != nil {
		return ""
	}

	return short
}

func (s *DBStorage) GetFullURL(shortURL string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ret string

	sql := "SELECT \"full\" FROM \"url\" where short = $1"
	row := s.database.QueryRowContext(ctx, sql, shortURL)

	err := row.Scan(&ret)

	if err != nil {
		fmt.Printf("\tNo URL in storage: %s\n", shortURL)
		return "", errors.New("wrong id")
	}
	return ret, nil
}

func (s *DBStorage) GetShortURL(fullURL string) (*ShortURL, error) {
	return &ShortURL{
		Short: "dummy",
	}, nil
}

func (s *DBStorage) GetCount() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ret int
	sql := "SELECT count(id) FROM \"url\""
	row := s.database.QueryRowContext(ctx, sql)
	err := row.Scan(&ret)

	if err != nil {
		return 0
	}
	return ret
}

func (s *DBStorage) GetUserURLs(user string) (string, error) {
	return "dummy", nil
}

func (s *DBStorage) AddUser() (string, string) {
	return s.userStorage.Add()
}

func (s *DBStorage) GetUser(hash string) string {
	return s.userStorage.GetUser(hash)
}

func (s *DBStorage) Load(val string) error {
	fmt.Printf("\tload function not implemented. file= %s not used\n", val)
	return nil
}
