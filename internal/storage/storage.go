package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"os"
	"strconv"
	"sync"
)

type URL struct {
	User  string `json:"user"`
	Full  string `json:"url"`
	Short string `json:"result"`
}

type ShortURL struct {
	Short string `json:"result"`
}

type FullURL struct {
	Full string `json:"url"`
}

type UserURL struct {
	Full  string `json:"original_url"`
	Short string `json:"short_url"`
}

type Storage interface {
	InsertURL(fullURL string) string
	GetFullURL(shortURL string) (string, error)
	GetShortURL(fullURL string) (*ShortURL, error)
	Load(fileName string) error
	GetCount() int
	GetUserURLs(user string) (string, error)
}

type MapStorage struct {
	counter         int
	Data            map[int]URL
	fileStorageName string
	mtx             sync.RWMutex
}

func New() Storage {
	return &MapStorage{
		counter: 0,
		Data:    make(map[int]URL),
	}
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

func (s *MapStorage) Load(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	s.fileStorageName = fileName
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
				fmt.Printf("\tLoad URL to storage = %s from file %s\n", URLItem, fileName)
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
				Short: element.Short,
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
