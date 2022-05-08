package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Sur0vy/url_shortner.git/internal/config"
	"github.com/Sur0vy/url_shortner.git/internal/database"
	"github.com/Sur0vy/url_shortner.git/internal/storage"
	"github.com/gin-gonic/gin"
	"io"
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

func NewBaseHandler(s *storage.Storage) *BaseHandler {
	return &BaseHandler{
		storage: *s,
	}
}

func (h *BaseHandler) GetFullURL(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	shortURL := c.Param("id")
	fmt.Printf("GetFullURL: short URL(param) = %s\n", shortURL)
	fullURL, err := h.storage.GetFullURL(shortURL)
	if err != nil {
		fmt.Printf("\tError: no full url")
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	} else {
		fmt.Printf("\tfull URL = %s\n", fullURL)
		if !strings.HasPrefix(fullURL, config.HTTP) {
			fullURL = config.HTTP + strings.TrimPrefix(fullURL, "//")
		}
		c.Redirect(http.StatusTemporaryRedirect, fullURL)
	}
}

func (h *BaseHandler) CreateShortURL(c *gin.Context) {
	var shortURL string
	fullURL, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("CreateShortURL: full URL(body) = %s\n", string(fullURL))
	shortURL = h.storage.InsertURL(string(fullURL))
	exShortURL := storage.ExpandShortURL(shortURL)
	fmt.Printf("\tShort URL = %s\n", exShortURL)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(http.StatusCreated, exShortURL)
}

func (h *BaseHandler) GetShortURL(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	fmt.Printf("GetShortURL: body(body) = %s\n", body)
	var fullURL storage.FullURL
	err = json.Unmarshal(body, &fullURL)
	if err != nil {
		fmt.Printf("\tError: no full URL")
		c.String(http.StatusBadRequest, "")
		return
	}

	fmt.Printf("\tfull URL = %s\n", fullURL.Full)
	var shortURL *storage.ShortURL
	var resStatus int
	shortURL, err = h.storage.GetShortURL(fullURL.Full)
	if err != nil {
		fmt.Printf("\tError: no short URL, creating\n")
		strURL := h.storage.InsertURL(fullURL.Full)
		shortURL = &storage.ShortURL{
			Short: strURL,
		}
		resStatus = http.StatusCreated
	} else {
		resStatus = http.StatusOK
	}
	shortURL.Short = storage.ExpandShortURL(shortURL.Short)
	fmt.Printf("\tshort URL = %s\n", shortURL.Short)
	body, err = json.Marshal(shortURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	}
	c.String(resStatus, string(body))
}

func (h *BaseHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
}

func (h *BaseHandler) GetUserURLs(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	users, err := h.storage.GetUserURLs(config.Cnf.CurrentUser)
	if err != nil {
		c.AbortWithStatus(http.StatusNoContent)
	} else {
		c.String(http.StatusOK, users)
	}
}

func (h *BaseHandler) Ping(c *gin.Context) {
	err := database.DB.Ping()
	if err == nil {
		c.Status(http.StatusOK)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *BaseHandler) AppendGroup(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	URLs, err := readURLsFromJSON(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	//сохранить URL в БД и сразу получаем JSON с сокращенными URL
	res, err := h.storage.InsertURLs(URLs)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.String(http.StatusCreated, res)
}

func readURLsFromJSON(reader io.ReadCloser) ([]storage.URLIdFull, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	fmt.Printf("GetShortURL: body(body) = %s\n", input)

	var URLs []storage.URLIdFull
	if err := json.Unmarshal(input, &URLs); err != nil {
		fmt.Printf("\tURL unmarshal error: %s\n", err)
		return nil, err
	}
	fmt.Printf("\tURL unmarshal success")
	return URLs, nil
}
