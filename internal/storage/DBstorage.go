package storage

import (
	"database/sql"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBStorage struct {
	connectionStr string
	database      *sql.DB
}

func NewDBStorage() Storage {
	s := &DBStorage{}
	err := s.Load(config.Cnf.DatabaseCon)
	if err != nil {
		fmt.Printf("\tNo stotage, or storage is corrapted!\n")
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

func (s *DBStorage) Load(val string) error {
	var err error
	s.connectionStr = val

	s.database, err = sql.Open("pgx", s.connectionStr)
	if err != nil {
		return err
	}
	defer s.database.Close()
	return nil
}

func (s *DBStorage) GetCount() int {
	return -1
}

func (s *DBStorage) GetUserURLs(user string) (string, error) {
	return "dummy", nil
}

func (s *DBStorage) IsAvailable() bool {
	return s.database.Ping() == nil
}
