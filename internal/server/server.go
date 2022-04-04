package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/handlers"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	memoryStorage := storage.NewMapStorage()
	handler := handlers.NewHandler(memoryStorage)

	router := gin.Default()
	router.GET("/:id", handler.GetURL)
	router.POST("/", handler.CreateShortURL)
	router.NoRoute(handler.ResponseError)
	return router
}
