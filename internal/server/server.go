package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/handlers"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	Users = NewMapUserStorage()
	Users.LoadFromFile()
	memoryStorage := storage.NewMapStorage()
	handler := handlers.NewBaseHandler(memoryStorage)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CookieMidlewared())
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.GET("/:id", handler.GetFullURL)
	router.GET("/api/user/urls", handler.GetUserURLs)
	router.POST("/", handler.CreateShortURL)
	router.POST("/api/shorten", handler.GetShortURL)
	router.GET("/ping", handler.Ping)
	router.NoRoute(handler.ResponseBadRequest)
	return router
}
