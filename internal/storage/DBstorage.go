package storage

//TODO иесиы на модуль

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/storage/helpers"
	"github.com/Sur0vy/url_shortner.git/internal/storage/users"
)

type DBStorage struct {
	userStorage users.UserStorage
	db          *sql.DB
	idle        *helpers.AtomicBool
}

func NewDBStorage(ctx context.Context) Storage {
	s := &DBStorage{
		idle: helpers.NewAtomicBool(true),
	}
	s.Connect()
	s.CreateTables(ctx)
	s.userStorage = users.NewDBUserStorage(s.db)
	return s
}

func (s *DBStorage) InsertURL(ctx context.Context, fullURL string) (string, error) {
	short := s.URLExist(ctx, fullURL)
	if short != "" {
		return short, NewURLError("URL is exist")
	}

	short = strconv.Itoa(s.GetCount(ctx) + 1)

	fmt.Printf("\tAdd new URL to storage full = %s, short = %s\n", fullURL, short)

	ctxIn, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	//добавление
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		helpers.TURL, helpers.FFull, helpers.FShort)
	_, err := s.db.ExecContext(ctxIn, sqlStr, fullURL, short)
	if err != nil {
		return "", err
	}

	//вставить ссылку на юзера
	sqlStr = fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		helpers.TUserURL, helpers.FShort, helpers.FUserHash)
	_, err = s.db.ExecContext(ctxIn, sqlStr, short, config.Cnf.CurrentUserHash)
	if err != nil {
		return "", err
	}

	return short, err
}

func (s *DBStorage) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var (
		URL    string
		status string
	)

	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s where %s = $1",
		helpers.FFull, helpers.FStatus, helpers.TURL, helpers.FShort)
	row := s.db.QueryRowContext(ctxIn, sqlStr, shortURL)
	err := row.Scan(&URL, &status)
	if err != nil {
		fmt.Printf("\tNo URL in storage: %s\n", shortURL)
		return "", errors.New("wrong id")
	} else if status == "del" {
		fmt.Printf("\tURL is gone in storage: %s\n", shortURL)
		return "", NewURLGoneError("URL is gone")
	}
	return URL, nil
}

func (s *DBStorage) GetShortURL(ctx context.Context, fullURL string) (*ShortURL, error) {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var short string

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		helpers.FShort, helpers.TURL, helpers.FFull)
	row := s.db.QueryRowContext(ctxIn, sqlStr, fullURL)

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

func (s *DBStorage) GetCount(ctx context.Context) int {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var ret int
	sqlStr := fmt.Sprintf("SELECT count(%s) FROM %s",
		helpers.FShort, helpers.TURL)
	row := s.db.QueryRowContext(ctxIn, sqlStr)
	err := row.Scan(&ret)

	if err != nil {
		return 0
	}
	return ret
}

func (s *DBStorage) GetUserURLs(ctx context.Context, user string) (string, error) {
	var userDataList []UserURL
	fmt.Printf("\tGet user %s urls\n", user)

	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlStr := fmt.Sprintf("SELECT u.%s, u.%s FROM %s u JOIN %s d ON u.%s = d.%s"+
		" JOIN %s r ON d.%s = r.%s WHERE r.%s = $1",
		helpers.FFull, helpers.FShort, helpers.TURL,
		helpers.TUserURL, helpers.FShort, helpers.FShort,
		helpers.TUser, helpers.FUserHash, helpers.FUserHash, helpers.FUserName)
	rows, err := s.db.QueryContext(ctxIn, sqlStr, user)
	if err != nil {
		return "", err
	}

	// обязательно закрываем перед возвратом функции
	defer func() {
		_ = rows.Close()
	}()

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

func (s *DBStorage) AddUser(ctx context.Context) (string, string) {
	return s.userStorage.Add(ctx)
}

func (s *DBStorage) GetUser(ctx context.Context, hash string) string {
	return s.userStorage.GetUser(ctx, hash)
}

func (s *DBStorage) Load(val string) error {
	fmt.Printf("\tload function not implemented. file= %s not used\n", val)
	return nil
}

func (s *DBStorage) URLExist(ctx context.Context, fullURL string) string {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var short string

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1",
		helpers.FShort, helpers.TURL, helpers.FFull)
	row := s.db.QueryRowContext(ctxIn, sqlStr, fullURL)

	err := row.Scan(&short)

	if err != nil {
		fmt.Printf("\tSQL count error: %s\n", fullURL)
		return ""
	}
	return short
}

func (s *DBStorage) InsertURLs(ctx context.Context, URLs []URLIdFull) (string, error) {
	ctxIn, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES($1, $2, $3) ON CONFLICT (%s) DO NOTHING",
		helpers.TURL, helpers.FInfo, helpers.FFull, helpers.FShort, helpers.FShort)
	insertStmt, err := s.db.PrepareContext(ctxIn, sqlStr)
	if err != nil {
		return "", err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	txStmt := tx.StmtContext(ctxIn, insertStmt)

	inc := 1
	existCnt := 0
	for i, val := range URLs {
		short := s.URLExist(ctxIn, URLs[i].Full)
		if short != "" {
			URLs[i].Short = short
			existCnt++
		} else {
			URLs[i].Short = strconv.Itoa(s.GetCount(ctx) + inc)
			inc++
			if _, err = txStmt.ExecContext(ctxIn, val.ID, val.Full, URLs[i].Short); err != nil {
				return "", err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	//встака ссылки на пользователя
	sqlStr = fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES($1, $2) ON CONFLICT (%s, %s) DO NOTHING",
		helpers.TUserURL, helpers.FShort, helpers.FUserHash, helpers.FShort, helpers.FUserHash)
	insertStmt, err = s.db.PrepareContext(ctxIn, sqlStr)
	if err != nil {
		return "", err
	}

	tx, err = s.db.Begin()
	if err != nil {
		return "", err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	txStmt = tx.StmtContext(ctxIn, insertStmt)

	for i, val := range URLs {
		if _, err = txStmt.ExecContext(ctxIn, val.Short, config.Cnf.CurrentUserHash); err != nil {
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
	if len(URLs) == existCnt {
		return string(data), NewURLError("URL is exist")
	}
	return string(data), nil
}

func (s *DBStorage) Connect() {
	var err error
	s.db, err = sql.Open("pgx", config.Cnf.DatabaseDSN)
	if err != nil {
		panic(err)
	}
}

func (s *DBStorage) CreateTables(ctx context.Context) {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT, %s TEXT UNIQUE, %s TEXT UNIQUE PRIMARY KEY, %s TEXT)",
		helpers.TURL, helpers.FInfo, helpers.FFull, helpers.FShort, helpers.FStatus)
	_, err := s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		panic(err)
	}

	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT NOT NULL, %s TEXT NOT NULL PRIMARY KEY)",
		helpers.TUser, helpers.FUserName, helpers.FUserHash)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		panic(err)
	}

	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT NOT NULL, %s TEXT NOT NULL, PRIMARY KEY (%s, %s))",
		helpers.TUserURL, helpers.FShort, helpers.FUserHash, helpers.FShort, helpers.FUserHash)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		panic(err)
	}

	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_u_us_URLID;"+
		"ALTER TABLE %s ADD CONSTRAINT fk_u_us_URLID FOREIGN KEY (%s) REFERENCES %s (%s);",
		helpers.TUserURL, helpers.TUserURL, helpers.FShort, helpers.TURL, helpers.FShort)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		panic(err)
	}

	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_u_us_UserID;"+
		"ALTER TABLE %s ADD CONSTRAINT fk_u_us_UserID FOREIGN KEY (%s) REFERENCES %s (%s);",
		helpers.TUserURL, helpers.TUserURL, helpers.FUserHash, helpers.TUser, helpers.FUserHash)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		panic(err)
	}
}

func (s *DBStorage) Ping() error {
	return s.db.Ping()
}

func (s *DBStorage) DeleteShortURLs(ctx context.Context, hash string, IDs []string) error {
	go func() {
		s.idle.Set(false)

		ctxIn, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		//добавление
		sqlStr := fmt.Sprintf("UPDATE %s AS u SET %s = $1 FROM %s AS d, %s AS r"+
			" WHERE (d.%s = u.%s) AND (r.%s = d.%s) AND"+
			" u.%s = ANY($2) AND (r.%s = $3)",
			helpers.TURL, helpers.FStatus, helpers.TUserURL, helpers.TUser,
			helpers.FShort, helpers.FShort, helpers.FUserHash, helpers.FUserHash,
			helpers.FShort, helpers.FUserHash)

		//ids := &pgtype.Int4Array{}
		//ids.Set(IDs)

		_, err := s.db.ExecContext(ctxIn, sqlStr, helpers.VDel, IDs, hash)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		defer s.idle.Set(true)
	}()
	return nil
}
