package router

import (
	"github.com/gin-gonic/gin"
	"go-job/master/api"
)

func NewWebRouter(mdls []gin.HandlerFunc,
	jobApi *api.JobApi,
	jobRecordApi *api.JobRecordApi,
	nodeApi *api.NodeApi,
	userApi *api.UserApi,
	dashboardApi *api.DashboardApi) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	group := server.Group("/api/go-job")
	jobApi.RegisterRoutes(group)
	jobRecordApi.RegisterRoutes(group)
	nodeApi.RegisterRoutes(group)
	userApi.RegisterRoutes(group)
	dashboardApi.RegisterRoutes(group)
	return server
}
