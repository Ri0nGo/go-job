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
		repo.NewJobRecordRepo,
		repo.NewNodeRepo,
		repo.NewUserRepo,

		// service
		service.NewJobService,
		service.NewJobRecordService,
		service.NewNodeService,
		service.NewUserService,

		// api
		api.NewJobApi,
		api.NewJobRecordApi,
		api.NewNodeApi,
		api.NewUserApi,

		// web
		middleware.NewGinMiddlewares,
		router.NewWebRouter,
	)
	return gin.Default()
}
