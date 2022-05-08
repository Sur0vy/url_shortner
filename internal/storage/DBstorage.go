package storage

//TODO иесиы на модуль

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/database"
	"github.com/Sur0vy/url_shortner.git/internal/storage/users"
	_ "github.com/jackc/pgx/v4/stdlib"
	"strconv"
	"time"
)

type DBStorage struct {
	userStorage users.UserStorage
	database    *sql.DB
}

func NewDBStorage(db *sql.DB) Storage {
	s := &DBStorage{
		database:    db,
		userStorage: users.NewDBUserStorage(db),
	}
	return s
}

func (s *DBStorage) InsertURL(fullURL string) string {
	//проверка на существование
	short := s.URLExist(fullURL)
	if short != "" {
		return short
	}

	short = strconv.Itoa(s.GetCount() + 1)

	fmt.Printf("\tAdd new URL to storage full = %s, short = %s\n", fullURL, short)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//добавление
	sql := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		database.TURL, database.FFull, database.FShort)
	_, err := s.database.ExecContext(ctx, sql, fullURL, short)
	if err != nil {
		return ""
	}

	//вставить ссылку на юзера
	sql = fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		database.TUserURL, database.FShort, database.FUserHash)
	_, err = s.database.ExecContext(ctx, sql, short, config.Cnf.CurrentUserHash)
	if err != nil {
		return ""
	}

	return short
}

func (s *DBStorage) GetFullURL(shortURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ret string

	sql := fmt.Sprintf("SELECT %s FROM %s where %s = $1",
		database.FFull, database.TURL, database.FShort)
	row := s.database.QueryRowContext(ctx, sql, shortURL)
	err := row.Scan(&ret)
	if err != nil {
		fmt.Printf("\tNo URL in storage: %s\n", shortURL)
		return "", errors.New("wrong id")
	}
	return ret, nil
}

func (s *DBStorage) GetShortURL(fullURL string) (*ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var short string

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		database.FShort, database.TURL, database.FFull)
	row := s.database.QueryRowContext(ctx, sql, fullURL)

	err := row.Scan(&short)

	if err != nil {
		fmt.Printf("\tNo URL in storage: %s\n", fullURL)
		return nil, errors.New("wrong URL")
	}
	fmt.Printf("\tGet short URL from storage. full = %s ; short = %s\n", fullURL, short)
	return &ShortURL{
		Short: short,
	}, nil
}

func (s *DBStorage) GetCount() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ret int
	sql := fmt.Sprintf("SELECT count(%s) FROM %s",
		database.FShort, database.TURL)
	row := s.database.QueryRowContext(ctx, sql)
	err := row.Scan(&ret)

	if err != nil {
		return 0
	}
	return ret
}

func (s *DBStorage) GetUserURLs(user string) (string, error) {
	var userDataList []UserURL
	fmt.Printf("\tGet user %s urls\n", user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := fmt.Sprintf("SELECT u.%s, u.%s FROM %s u JOIN %s d ON u.%s = d.%s"+
		" JOIN %s r ON d.%s = r.%s WHERE r.%s = $1",
		database.FFull, database.FShort, database.TURL,
		database.TUserURL, database.FShort, database.FShort,
		database.TUser, database.FUserHash, database.FUserHash, database.FUserName)
	rows, err := s.database.QueryContext(ctx, sql, user)
	if err != nil {
		return "", err
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var userData UserURL
		err = rows.Scan(&userData.Full, &userData.Short)
		if err != nil {
			return "", err
		}
		userDataList = append(userDataList, userData)
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	if len(userDataList) == 0 {
		return "", errors.New("no URLs found")
	}
	data, err := json.Marshal(&userDataList)
	if err != nil {
		return "", err
	}
	return string(data), nil
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

func (s *DBStorage) URLExist(fullURL string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var short string

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1",
		database.FShort, database.TURL, database.FFull)
	row := s.database.QueryRowContext(ctx, sql, fullURL)

	err := row.Scan(&short)

	if err != nil {
		fmt.Printf("\tSQL count error: %s\n", fullURL)
		return ""
	}
	return short
}

func (s *DBStorage) InsertURLs(URLs []URLIdFull) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sql := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES($1, $2, $3)",
		database.TURL, database.FInfo, database.FFull, database.FShort)
	insertStmt, err := s.database.PrepareContext(ctx, sql)

	tx, err := s.database.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, insertStmt)

	inc := 1
	for i, _ := range URLs {
		short := s.URLExist(URLs[i].Full)
		if short != "" {
			URLs[i].Short = short
		} else {
			URLs[i].Short = strconv.Itoa(s.GetCount() + inc)
			inc++
			if _, err = txStmt.ExecContext(ctx, URLs[i].Id, URLs[i].Full, URLs[i].Short); err != nil {
				return "", err

			}
		}
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	//встака ссылки на пользователя
	sql = fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES($1, $2)",
		database.TUserURL, database.FShort, database.FUserHash)
	insertStmt, err = s.database.PrepareContext(ctx, sql)

	tx, err = s.database.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	txStmt = tx.StmtContext(ctx, insertStmt)

	for i, _ := range URLs {
		if _, err = txStmt.ExecContext(ctx, URLs[i].Short, config.Cnf.CurrentUserHash); err != nil {
			return "", err
		}
		URLs[i].Full = ""
		URLs[i].Short = ExpandShortURL(URLs[i].Short)
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	//json to string marshall
	data, err := json.Marshal(&URLs)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
