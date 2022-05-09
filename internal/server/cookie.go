package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func CookieMidlewared(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("url_shortner")
		if err == nil && cookie != "" {
			user := (*s).GetUser(cookie)
			if user != "" {
				config.Cnf.CurrentUser = user
				config.Cnf.CurrentUserHash = cookie
				return
			}
		}
		user, hash := (*s).AddUser()
		config.Cnf.CurrentUser = user
		config.Cnf.CurrentUserHash = hash
		c.SetCookie("url_shortner", hash, 3600, "/", "localhost", false, true)
	}
}
