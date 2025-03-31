package router

import (
	"github.com/gin-gonic/gin"
	"go-job/node/api"
)

func RegistryRoute(engine *gin.Engine) {
	apiG := engine.Group("/api/")

	apiG.POST("/job/add", api.AddJob)
}
