package app

import (
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetURL(c *gin.Context) {
	shortURL := c.Param("id")
	fullURL, err := h.storage.Get(shortURL)
	if err != nil {
		c.String(404, "")
	} else {
		c.Status(307)
		c.Header("Location", fullURL)
	}
}

func (h *Handler) CreateShortURL(c *gin.Context) {
	fullURL := c.Param("url")
	shortURL := h.storage.Insert(fullURL)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(201, shortURL)
}

func (h *Handler) ResponseError(c *gin.Context) {
	c.String(400, "")
}
