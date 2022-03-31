package server

import (
	"github.com/Sur0vy/url_shortner.git/internal/app"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
	"strconv"
)

func StartServer(port int) {
	memoryStorage := storage.NewMapStorage()
	handler := app.NewHandler(memoryStorage)

	router := gin.Default()
	router.GET("/:id", handler.GetURL)
	router.POST("/", handler.CreateShortURL)
	router.Run(generateAddress(port))
}

func generateAddress(port int) string {
	return ":" + strconv.Itoa(port)
}