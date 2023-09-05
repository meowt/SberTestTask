package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func SetupProbesHandler() Handler {
	return Handler{}
}

func (h *Handler) InitProbesRoutes(router *gin.Engine) {
	probesRouter := router.Group("/probes")
	{
		probesRouter.GET("/liveness", h.Liveness)
		probesRouter.GET("/readiness", h.Readiness)
	}
}

func (h *Handler) Liveness(c *gin.Context) {
	//checking health of app
	c.Data(http.StatusOK, "text/plain", []byte("OK"))
}

func (h *Handler) Readiness(c *gin.Context) {
	//checking is the app ready
	c.Data(http.StatusOK, "text/plain", []byte("OK"))
}
