//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-job/internal/pkg/email"
	"go-job/master/api"
	"go-job/master/database"
	"go-job/master/pkg/middleware"
	"go-job/master/repo"
	"go-job/master/repo/cache"
	"go-job/master/router"
	"go-job/master/service"
	"gorm.io/gorm"
)

type WebContainer struct {
	Engine  *gin.Engine
	MysqlDB *gorm.DB
	JobSvc  service.IJobService
}

func InitWebServer() *WebContainer {
	wire.Build(
		//database
		database.NewMySQLWithGORM,
		database.NewRedisClient,
		cache.NewEmailCodeCache,

		// repo
		repo.NewJobRepo,
		repo.NewJobRecordRepo,
		repo.NewNodeRepo,
		repo.NewUserRepo,
		repo.NewEmailCodeRepo,

		// service
		service.NewJobService,
		service.NewJobRecordService,
		service.NewNodeService,
		email.InitEmailService,
		service.NewEmailCodeService,
		service.NewUserService,

		// api
		api.NewJobApi,
		api.NewJobRecordApi,
		api.NewNodeApi,
		api.NewUserApi,

		// web
		middleware.NewGinMiddlewares,
		router.NewWebRouter,

		wire.Struct(new(WebContainer), "*"),
	)
	return nil
}
