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
	c.String(200, "Success")
}

func (h *Handler) CreateShortURL(c *gin.Context) {
	c.String(200, "qqqqq")
}
