package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/gin-gonic/gin"
)

var Users UserStorage

func CookieMidlewared() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("url_shortner")
		if err != nil || cookie == "" {
			user, id := Users.Add()
			config.Cnf.CurrentUser = user
			c.SetCookie("url_shortner", id, 3600, "/", "localhost", false, true)
		} else {
			config.Cnf.CurrentUser = Users.GetUser(cookie)
		}
	}
}
