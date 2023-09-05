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
	PutObject(key, value string, ttl time.Duration) (res string, err error)
	GetObject(key string) (obj models.Object, err error)
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

	res, err := h.Repo.PutObject(key, string(jsonData), ttl)
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetObject(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, "empty \"key\" parameter")
		return
	}

	obj, err := h.Repo.GetObject(key)
	if obj.Data == "" {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	c.Data(http.StatusOK, "application-json", []byte(obj.Data))
}
