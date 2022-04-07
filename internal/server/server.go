package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/handlers"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	memoryStorage := storage.NewMapStorage()
	handler := handlers.NewBaseHandler(memoryStorage)

	router := gin.Default()
	router.GET("/:id", handler.GetFullURL)
	router.POST("/", handler.CreateShortURL)
	router.POST("/api/shorten", handler.GetShortURL)
	router.NoRoute(handler.ResponseBadRequest)
	return router
}
