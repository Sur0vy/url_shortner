package storage

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

type URLIdFull struct {
	ID    string `json:"correlation_id,omitempty"`
	Full  string `json:"original_url,omitempty"`
	Short string `json:"short_url,omitempty"`
}

type Storage interface {
	InsertURL(fullURL string) (string, error)
	GetFullURL(shortURL string) (string, error)
	GetShortURL(fullURL string) (*ShortURL, error)
	Load(val string) error
	GetCount() int
	GetUserURLs(user string) (string, error)
	AddUser() (string, string)
	GetUser(hash string) string
	InsertURLs([]URLIdFull) (string, error)
}
