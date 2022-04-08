package handlers

import (
	"encoding/json"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
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
		//c.Status(200)
		if !strings.HasPrefix(fullURL, config.Params.BaseURL) {
			fullURL = config.Params.BaseURL + strings.TrimPrefix(fullURL, "//")
		}
		c.Writer.Header().Set("Location", fullURL)
	}
}

func (h *BaseHandler) CreateShortURL(c *gin.Context) {
	var shortURL string
	fullURL, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		shortURL = ""
	}
	shortURL = config.LocalHost + h.storage.InsertURL(string(fullURL))
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
	var resStatus int
	shortURL, err = h.storage.GetShortURL(fullURL.Full)
	if err != nil {
		//тогда создадим
		strURL := h.storage.InsertURL(string(fullURL.Full))
		shortURL = &storage.ShortURL{
			Short: strURL,
		}
		resStatus = http.StatusCreated
	} else {
		resStatus = http.StatusOK
	}
	body, err = json.Marshal(shortURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	}
	c.String(resStatus, string(body))
}

func (h *BaseHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
}
