package storage

import (
	"errors"
	"strconv"
	"sync"
)

type URL struct {
	Full  string `json:"full"`
	Short string `json:"short"`
}
type Storage interface {
	Insert(fullURL string) string
	Get(id string) (string, error)
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

func (s *MapStorage) Insert(fullURL string) string {
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

func (s *MapStorage) Get(shortURL string) (string, error) {
	//пока код не имеет значения
	for _, element := range s.Data {
		if element.Short == shortURL {
			return element.Full, nil
		}
	}
	return "", errors.New("wrong id")
}
