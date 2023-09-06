package handler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/meowt/SberTestTask/pkg/models/objects"
)

type RepoService interface {
	PutObject(key, value string, ttl time.Duration)
	GetObject(key string) (obj models.Object)
}

type Handler struct {
	Repo RepoService
}

func SetupObjectsHandler(r RepoService) Handler {
	return Handler{Repo: r}
}

func (h *Handler) InitObjectsRoutes(router *gin.Engine) {
	objectsRouter := router.Group("/objects")
	{
		objectsRouter.PUT("/:key", h.PutObject)
		objectsRouter.GET("/:key", h.GetObject)
	}
}

func (h *Handler) PutObject(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		err = c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}
	key := c.Param("key")
	if key == "" {
		err = c.AbortWithError(http.StatusBadRequest, errors.New("empty \"key\" parameter"))
		log.Println(err)
		return
	}

	var ttl time.Duration
	expires := c.GetHeader("Expires")
	if expires != "" {
		exp, err := strconv.Atoi(expires)
		if err != nil {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}
		ttl = time.Duration(exp) * time.Millisecond
	}

	h.Repo.PutObject(key, string(jsonData), ttl)
	c.JSON(http.StatusOK, "OK")
}

func (h *Handler) GetObject(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, "empty \"key\" parameter")
		return
	}

	obj := h.Repo.GetObject(key)
	if obj.Data == "" {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.Data(http.StatusOK, "application-json", []byte(obj.Data))
}
