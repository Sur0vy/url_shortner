package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/storage/users"
	"os"
	"strconv"
	"sync"
)

type MapStorage struct {
	counter         int
	Data            map[int]URL
	fileStorageName string
	userStorage     users.UserStorage
	mtx             sync.RWMutex
}

func NewMapStorage() Storage {
	s := &MapStorage{
		counter:     0,
		Data:        make(map[int]URL),
		userStorage: users.NewMapUserStorage(),
	}

	_ = s.userStorage.LoadFromFile()

	err := s.Load(config.Cnf.StoragePath)
	if err != nil {
		fmt.Printf("\tNo stotage, or storage is corrapted!\n")
	}
	return s
}

func (s *MapStorage) InsertURL(ctx context.Context, fullURL string) (string, error) {
	short, err := s.GetShortURL(ctx, fullURL)
	if err == nil {
		return short.Short, NewURLError("URL is exist")
	}
	s.counter++
	var URLItem = URL{
		User:  config.Cnf.CurrentUser,
		Full:  fullURL,
		Short: strconv.Itoa(s.counter),
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	fmt.Printf("\tAdd new URL to storage = %s\n", URLItem)
	s.Data[s.counter] = URLItem
	if len(s.fileStorageName) > 0 {
		err = s.addToFile(&URLItem)
		if err != nil {
			return "", err
		}
	}
	return URLItem.Short, nil
}

func (s *MapStorage) GetFullURL(_ context.Context, shortURL string) (string, error) {
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

func (s *MapStorage) GetShortURL(_ context.Context, fullURL string) (*ShortURL, error) {
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
	defer func() {
		_ = file.Close()
	}()

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
	defer func() {
		_ = file.Close()
	}()

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

func (s *MapStorage) GetCount(_ context.Context) int {
	return s.counter
}

func (s *MapStorage) GetUserURLs(_ context.Context, user string) (string, error) {
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

func (s *MapStorage) AddUser(ctx context.Context) (string, string) {
	return s.userStorage.Add(ctx)
}

func (s *MapStorage) GetUser(ctx context.Context, hash string) string {
	return s.userStorage.GetUser(ctx, hash)
}

func (s *MapStorage) InsertURLs(_ context.Context, _ []URLIdFull) (string, error) {
	//TODO реализовать логику + проверка на попытку вставить дубликат
	return "", nil
}

func (s *MapStorage) InvokeDeferFunction() {
	//TODO тут нечего делать, просто заглушка
}

func (s *MapStorage) Ping() error {
	return nil
}
