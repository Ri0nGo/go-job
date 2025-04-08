package router

import (
	"github.com/gin-gonic/gin"
	"go-job/master/api"
)

func NewWebRouter(mdls []gin.HandlerFunc,
	jobApi *api.JobApi,
	nodeApi *api.NodeApi) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	group := server.Group("/api/go-job")
	jobApi.RegisterRoutes(group)
	nodeApi.RegisterRoutes(group)
	return server
}
