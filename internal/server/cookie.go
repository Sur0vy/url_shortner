package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Users UserStorage

func CookieMidlewared() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("url_shortner")
		if err != nil || cookie == "" {
			c.String(http.StatusNoContent, "")
			user, id := Users.Add()
			config.Cnf.CurrentUser = user
			c.SetCookie("url_shortner", id, 3600, "/", "localhost", false, true)
		} else {
			config.Cnf.CurrentUser = Users.GetUser(cookie)
		}
	}
}
