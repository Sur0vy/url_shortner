package app

import (
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
		c.Writer.Header().Set("Location", fullURL)
	}
}

func (h *Handler) CreateShortURL(c *gin.Context) {
	var shortURL string
	fullURL, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		shortURL = ""
	}
	shortURL = "http://localhost:8080/" + h.storage.Insert(string(fullURL))
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(201, shortURL)
}

func (h *Handler) ResponseError(c *gin.Context) {
	c.String(400, "")
}
