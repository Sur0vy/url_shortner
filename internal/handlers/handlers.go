package handlers

import (
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type Handle interface {
	GetURL(c *gin.Context)
	CreateShortURL(c *gin.Context)
	ResponseBadRequest(c *gin.Context)
}

type BaseHandler struct {
	storage storage.Storage
}

func NewBaseHandler(storage storage.Storage) *BaseHandler {
	return &BaseHandler{storage: storage}
}

func (h *BaseHandler) GetURL(c *gin.Context) {
	shortURL := c.Param("id")
	fullURL, err := h.storage.Get(shortURL)
	if err != nil {
		c.String(http.StatusNotFound, "")
	} else {
		c.Status(http.StatusTemporaryRedirect)
		c.Writer.Header().Set("Location", fullURL)
	}
}

func (h *BaseHandler) CreateShortURL(c *gin.Context) {
	var shortURL string
	fullURL, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		shortURL = ""
	}
	shortURL = "http://localhost:8080/" + h.storage.Insert(string(fullURL))
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
}

func (h *BaseHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
}
