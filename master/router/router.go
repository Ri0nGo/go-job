package router

import (
	"github.com/gin-gonic/gin"
	"go-job/master/api"
)

func NewWebRouter(mdls []gin.HandlerFunc, userApi *api.JobApi) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	group := server.Group("/api/go-job")
	userApi.RegisterRoutes(group)
	return server
}
