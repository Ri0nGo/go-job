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
	dashboardApi *api.DashboardApi,
	iamOAuthApi *api.IAMOAuthApi,
	oauth2Api *api.OAuth2Api) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	group := server.Group("/api/go-job")
	jobApi.RegisterRoutes(group)
	jobRecordApi.RegisterRoutes(group)
	nodeApi.RegisterRoutes(group)
	userApi.RegisterRoutes(group)
	dashboardApi.RegisterRoutes(group)
	iamOAuthApi.RegisterRoutes(group)
	// oauth2Api.RegisterRoutes(group)
	return server
}
