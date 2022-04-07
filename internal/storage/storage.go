package storage

import (
	"errors"
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
		Short: strconv.Itoa(s.Counter),
	}
	s.Data[s.Counter] = url
	s.Unlock()
	return url.Short
}

func (s *MapStorage) GetFullURL(shortURL string) (string, error) {
	//пока код не имеет значения
	for _, element := range s.Data {
		if element.Short == shortURL {
			return element.Full, nil
		}
	}
	return "", errors.New("wrong id")
}

func (s *MapStorage) GetShortURL(fullURL string) (*ShortURL, error) {
	//пока код не имеет значения
	for _, element := range s.Data {
		if element.Full == fullURL {
			return &ShortURL{element.Short}, nil
		}
	}
	return nil, errors.New("wrong URL")
}
