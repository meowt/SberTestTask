package router

import (
	"github.com/gin-gonic/gin"
	"github.com/meowt/SberTestTask/cmd/modules"
)

type Router struct {
	*gin.Engine
}

func Setup(h *modules.HandlerModule) (r *Router) {
	r = &Router{gin.New()}
	r.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	r.InitRoutes(h)
	return
}

func (r *Router) InitRoutes(h *modules.HandlerModule) {
	h.MetricsHandler.InitMetricsRoutes(r.Engine)
	h.ObjectsHandler.InitObjectsRoutes(r.Engine)
	h.ProbesHandler.InitProbesRoutes(r.Engine)
}
