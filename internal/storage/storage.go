package storage

import (
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
	return s.getFullURL(url), nil
}

func (s *MapStorage) getFullURL(shortUrl string) string {
	//пока код не имеет значения
	for _, element := range s.data {
		// element is the element from someSlice for where we are
		if element.Short == shortUrl {
			return element.Full
		}
	}
	return ""
}
