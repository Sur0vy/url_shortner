package storage

import (
	"errors"
	"strconv"
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
	counter int
	data    map[int]URL
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		counter: 0,
		data:    make(map[int]URL),
	}
}

func (s *MapStorage) Insert(fullURL string) string {
	//пока код не имеет значения
	s.counter++
	var url = URL{
		Full:  fullURL,
		Short: strconv.Itoa(s.counter),
	}
	s.data[s.counter] = url
	return url.Short
}

func (s *MapStorage) Get(url string) (string, error) {
	//пока код не имеет значения
	fullURL, err := s.getFullURL(url)
	if err != nil {
		return "", errors.New("wrong id")
	}
	return fullURL, nil
}

func (s *MapStorage) getFullURL(shortURL string) (string, error) {
	//пока код не имеет значения
	for _, element := range s.data {
		// element is the element from someSlice for where we are
		if element.Short == shortURL {
			return element.Full, nil
		}
	}
	return "", errors.New("wrong id")
}
