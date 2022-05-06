package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/storage/users"
	_ "github.com/jackc/pgx/v4/stdlib"
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
	return "dummy"
}

func (s *DBStorage) GetFullURL(shortURL string) (string, error) {
	return "dummy", nil
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
