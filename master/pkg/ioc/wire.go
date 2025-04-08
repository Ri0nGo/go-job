//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-job/master/api"
	"go-job/master/database"
	"go-job/master/pkg/middleware"
	"go-job/master/repo"
	"go-job/master/router"
	"go-job/master/service"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//database
		database.NewMySQLWithGORM,

		// repo
		repo.NewJobRepo,
		repo.NewNodeRepo,

		// service
		service.NewJobService,
		service.NewNodeService,

		// api
		api.NewJobApi,
		api.NewNodeApi,

		// web
		middleware.NewGinMiddlewares,
		router.NewWebRouter,
	)
	return gin.Default()
}
