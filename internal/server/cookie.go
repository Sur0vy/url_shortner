package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func CookieMidlewared(s *storage.Storage) gin.HandlerFunc {
	//TODO если пришел с кукой запрос, а ее нет в БД, нужно это обрабатывать
	return func(c *gin.Context) {
		cookie, err := c.Cookie("url_shortner")
		if err != nil || cookie == "" {
			user, hash := (*s).AddUser()
			config.Cnf.CurrentUser = user
			config.Cnf.CurrentUserHash = hash
			c.SetCookie("url_shortner", hash, 3600, "/", "localhost", false, true)
		} else {
			config.Cnf.CurrentUser = (*s).GetUser(cookie)
			config.Cnf.CurrentUserHash = cookie
		}
	}
}
