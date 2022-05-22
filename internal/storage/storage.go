package storage

import "context"

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
	InsertURL(ctx context.Context, fullURL string) (string, error)
	GetFullURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, fullURL string) (*ShortURL, error)
	Load(val string) error
	GetCount(ctx context.Context) int
	GetUserURLs(ctx context.Context, user string) (string, error)
	AddUser(ctx context.Context) (string, string)
	GetUser(ctx context.Context, hash string) string
	InsertURLs(ctx context.Context, URLs []URLIdFull) (string, error)
	DeleteShortURLs(ctx context.Context, hash string, IDs []string) error
	Ping() error
}
