package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/handlers"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupServer(s *storage.Storage) *gin.Engine {
	handler := handlers.NewBaseHandler(s)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CookieMidlewared(s))
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	api := router.Group("/api")
	api.GET("/user/urls", handler.GetUserURLs)
	api.DELETE("/user/urls", handler.DeleteUserURLs)
	api.POST("/shorten", handler.GetShortURL)
	api.POST("/shorten/batch", handler.AppendGroup)

	router.POST("/", handler.CreateShortURL)
	router.GET("/:id", handler.GetFullURL)
	router.GET("/ping", handler.Ping)

	router.NoRoute(handler.ResponseBadRequest)
	return router
}
