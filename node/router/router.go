package router

import (
	"github.com/gin-gonic/gin"
	"go-job/node/api"
	"go-job/node/service"
)

func InitRouter(engine *gin.Engine) {
	group := engine.Group("/api/go-job/node")

	jobService := service.NewJobService()
	jh := api.NewJobHandler(jobService)
	jh.RegisterRoutes(group)

	nodeService := service.NewNodeService()
	nh := api.NewNodeHandler(nodeService)
	nh.RegisterRoutes(group)
}
