package storage

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
)

func ExpandShortURL(shortURL string) string {
	return config.HTTP + config.HostAddr + ":" + config.Params.BasePort() + "/" + shortURL
}
