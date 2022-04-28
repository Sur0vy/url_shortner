package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"os"
	"strconv"
	"sync"
)

type MapStorage struct {
	counter         int
	Data            map[int]URL
	fileStorageName string
	mtx             sync.RWMutex
}

func NewMapStorage() Storage {
	s := &MapStorage{
		counter: 0,
		Data:    make(map[int]URL),
	}

	err := s.Load(config.Cnf.StoragePath)
	if err != nil {
		fmt.Printf("\tNo stotage, or storage is corrapted!\n")
	}
	return s
}

func (s *MapStorage) InsertURL(fullURL string) string {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.counter++
	var URLItem = URL{
		User:  config.Cnf.CurrentUser,
		Full:  fullURL,
		Short: strconv.Itoa(s.counter),
	}
	fmt.Printf("\tAdd new URL to storage = %s\n", URLItem)
	s.Data[s.counter] = URLItem
	if len(s.fileStorageName) > 0 {
		s.addToFile(&URLItem)
	}
	return URLItem.Short
}

func (s *MapStorage) GetFullURL(shortURL string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, element := range s.Data {
		if element.Short == shortURL {
			fmt.Printf("\tGet full URL from storage. short = %s ; full = %s\n", element.Short, element.Full)
			return element.Full, nil
		}
	}
	fmt.Printf("\tNo URL in storage: %s\n", shortURL)
	return "", errors.New("wrong id")
}

func (s *MapStorage) GetShortURL(fullURL string) (*ShortURL, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, element := range s.Data {
		if element.Full == fullURL {
			fmt.Printf("\tGet short URL from storage. full = %s ; short = %s\n", fullURL, element.Short)
			return &ShortURL{
				Short: element.Short,
			}, nil
		}
	}
	fmt.Printf("\tNo URL in storage: %s\n", fullURL)
	return nil, errors.New("wrong URL")
}

func (s *MapStorage) Load(val string) error {
	file, err := os.OpenFile(val, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	s.fileStorageName = val
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for {
		if !scanner.Scan() {
			return scanner.Err()
		}
		// читаем данные из scanner
		data := scanner.Bytes()
		URLItem := URL{}
		err = json.Unmarshal(data, &URLItem)
		if err == nil {
			func() {
				s.mtx.Lock()
				defer s.mtx.Unlock()
				s.counter++
				fmt.Printf("\tLoad URL to storage = %s from file %s\n", URLItem, val)
				s.Data[s.counter] = URLItem
			}()
		}
	}
}

func (s *MapStorage) addToFile(URLItem *URL) error {
	file, err := os.OpenFile(s.fileStorageName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(&URLItem)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	// записываем URL в буфер
	if _, err := writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return writer.Flush()
}

func (s *MapStorage) GetCount() int {
	return s.counter
}

func (s *MapStorage) GetUserURLs(user string) (string, error) {
	var userDataList []UserURL
	fmt.Printf("\tGet user %s urls\n", user)
	for _, element := range s.Data {
		if element.User == user {
			fmt.Printf("\t\tURL = %s\n", element.Full)
			userData := &UserURL{
				Full:  element.Full,
				Short: ExpandShortURL(element.Short),
			}
			userDataList = append(userDataList, *userData)
		}
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

func (s *MapStorage) IsAvailable() bool {
	//------------заглушка для прохождения тестов, пока БД не спроектирована и методы не описаны
	DSN := config.Cnf.DatabaseCon
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return false
	}
	return true
}
