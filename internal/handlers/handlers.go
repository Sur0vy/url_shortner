package handlers

import (
	"encoding/json"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type Handle interface {
	GetFullURL(c *gin.Context)
	CreateShortURL(c *gin.Context)
	GetShortURL(c *gin.Context)
	ResponseBadRequest(c *gin.Context)
}

type BaseHandler struct {
	storage storage.Storage
}

func NewBaseHandler(storage storage.Storage) *BaseHandler {
	return &BaseHandler{storage: storage}
}

func (h *BaseHandler) GetFullURL(c *gin.Context) {
	shortURL := c.Param("id")
	fullURL, err := h.storage.GetFullURL(shortURL)
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
	shortURL = "http://localhost:8080/" + h.storage.InsertURL(string(fullURL))
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
}

func (h *BaseHandler) GetShortURL(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	var fullURL storage.FullURL
	err = json.Unmarshal(body, &fullURL)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	var shortURL *storage.ShortURL
	shortURL, err = h.storage.GetShortURL(fullURL.Full)
	if err != nil {
		c.String(http.StatusNotFound, "")
	} else {
		body, err := json.Marshal(shortURL)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(body))
	}
}

func (h *BaseHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
}
