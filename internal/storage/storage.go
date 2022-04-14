package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type URL struct {
	Full  string `json:"url"`
	Short string `json:"result"`
}

type ShortURL struct {
	Short string `json:"result"`
}

type FullURL struct {
	Full string `json:"url"`
}

type Storage interface {
	InsertURL(fullURL string) string
	GetFullURL(shortURL string) (string, error)
	GetShortURL(fullURL string) (*ShortURL, error)
	Load(fileName string) error
}

type MapStorage struct {
	Counter         int
	Data            map[int]URL
	fileStorageName string
	sync.Mutex
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		Counter: 0,
		Data:    make(map[int]URL),
	}
}

func (s *MapStorage) InsertURL(fullURL string) string {
	//пока код не имеет значения
	s.Lock()
	s.Counter++
	var URLItem = URL{
		Full:  fullURL,
		Short: strconv.Itoa(s.Counter),
	}
	fmt.Printf("\tAdd new URL to storage = %s\n", URLItem)
	s.Data[s.Counter] = URLItem
	if len(s.fileStorageName) > 0 {
		s.addToFile(&URLItem)
	}
	s.Unlock()
	return URLItem.Short
}

func (s *MapStorage) GetFullURL(shortURL string) (string, error) {
	//пока код не имеет значения
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
	//пока код не имеет значения
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
			s.Lock()
			s.Counter++
			fmt.Printf("\tLoad URL to storage = %s from file %s\n", URLItem, fileName)
			s.Data[s.Counter] = URLItem
			s.Unlock()
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
