package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func CookieMidlewared(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("url_shortner")
		if err != nil || cookie == "" {
			user, id := (*s).AddUser()
			config.Cnf.CurrentUser = user
			c.SetCookie("url_shortner", id, 3600, "/", "localhost", false, true)
		} else {
			config.Cnf.CurrentUser = (*s).GetUser(cookie)
		}
	}
}
