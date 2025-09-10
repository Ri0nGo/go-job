package router

import (
	"github.com/gin-gonic/gin"
	"go-job/node/api"
)

func NewWebRouter(jobApi *api.JobApi, nodeApi *api.NodeApi) *gin.Engine {
	server := gin.Default()
	group := server.Group("/api/go-job/node")
	jobApi.RegisterRoutes(group)
	nodeApi.RegisterRoutes(group)
	return server
}
