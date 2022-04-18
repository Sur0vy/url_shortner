package server

import (
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/handlers"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	memoryStorage := storage.NewMapStorage()
	err := memoryStorage.Load(config.Cnf.StoragePath)
	if err != nil {
		fmt.Printf("\tNo stotage, or storage is corrapted!\n")
	}
	handler := handlers.NewBaseHandler(memoryStorage)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.GET("/:id", handler.GetFullURL)
	router.POST("/", handler.CreateShortURL)
	router.POST("/api/shorten", handler.GetShortURL)
	router.NoRoute(handler.ResponseBadRequest)
	return router
}
