//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-job/node/api"
	"go-job/node/router"
	"go-job/node/service"
)

type WebContainer struct {
	Engine *gin.Engine
	JobSvc service.IJobService
}

func InitWebServer() *WebContainer {
	wire.Build(

		// service
		service.NewJobService,
		service.NewNodeService,

		// api
		api.NewJobApi,
		api.NewNodeApi,

		router.NewWebRouter,

		wire.Struct(new(WebContainer), "*"),
	)
	return nil
}
