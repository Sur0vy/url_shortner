package storage

import "strconv"

type URL struct {
	Full  string `json:"full"`
	Short string `json:"short"`
}
type Storage interface {
	Insert(fullURL string)
	Get(id string) string
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

func (s *MapStorage) Insert(fullURL string) {
	//пока код не имеет значения
	s.counter++
	var url = URL{
		Full:  fullURL,
		Short: strconv.Itoa(s.counter),
	}
	s.data[s.counter] = url
}

func (s *MapStorage) Get(url string) string {
	//пока код не имеет значения
	return ""
}
