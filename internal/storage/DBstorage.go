package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBStorage struct {
	//	connectionStr string
	database *sql.DB
}

func NewDBStorage(db *sql.DB) Storage {
	s := &DBStorage{
		database: db,
	}
	//	err := s.Load(config.Cnf.DatabaseDSN)
	//	if err != nil {
	//		fmt.Printf("\tNo stotage, or storage is corrapted!\n")
	//	}
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

//func (s *DBStorage) Load(val string) error {
//	var err error
//	s.connectionStr = val
//
//	s.database, err = sql.Open("pgx", s.connectionStr)
//	if err != nil {
//		return err
//	}
//	defer s.database.Close()
//	return nil
//}

func (s *DBStorage) GetCount() int {
	return -1
}

func (s *DBStorage) GetUserURLs(user string) (string, error) {
	return "dummy", nil
}

func (s *DBStorage) AddUser() (string, string) {
	return "", ""
}

func (s *DBStorage) GetUser(hash string) string {
	return ""
}

func (s *DBStorage) Load(val string) error {
	return nil
}
