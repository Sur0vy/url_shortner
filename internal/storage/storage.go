package storage

import (
	"errors"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
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
}

type MapStorage struct {
	Counter int
	Data    map[int]URL
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
	var url = URL{
		Full:  fullURL,
		Short: config.HTTPPref + "/" + strconv.Itoa(s.Counter),
	}
	fmt.Printf("\t Add new URL to storage = %s\n", url)
	s.Data[s.Counter] = url
	s.Unlock()
	return url.Short
}

func (s *MapStorage) GetFullURL(shortURL string) (string, error) {
	//пока код не имеет значения
	for _, element := range s.Data {
		exShortURL := config.HTTPPref + "/" + shortURL
		if element.Short == exShortURL {
			fmt.Printf("\t Get full URL from storage. short = %s ; full = %s\n", exShortURL, element.Full)
			return element.Full, nil
		}
	}
	fmt.Printf("\t No URL in storage: %s\n", shortURL)
	return "", errors.New("wrong id")
}

func (s *MapStorage) GetShortURL(fullURL string) (*ShortURL, error) {
	//пока код не имеет значения
	for _, element := range s.Data {
		if element.Full == fullURL {
			fmt.Printf("\t Get short URL from storage. full = %s ; short = %s\n", fullURL, element.Short)
			return &ShortURL{
				Short: element.Short,
			}, nil
		}
	}
	fmt.Printf("\t No URL in storage: %s\n", fullURL)
	return nil, errors.New("wrong URL")
}
