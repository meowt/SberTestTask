package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
}

func SetupMetricsHandler() Handler {
	return Handler{}
}

func (h *Handler) InitMetricsRoutes(router *gin.Engine) {
	metricsRouter := router.Group("/metrics")
	{
		metricsRouter.GET("/", h.GetMetrics)
	}
}

func (h *Handler) GetMetrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
