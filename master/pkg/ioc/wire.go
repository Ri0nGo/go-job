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
	"go-job/master/pkg/notify"
	"go-job/master/repo"
	"go-job/master/repo/cache"
	"go-job/master/router"
	"go-job/master/service"
	"gorm.io/gorm"
)

type WebContainer struct {
	Engine      *gin.Engine
	MysqlDB     *gorm.DB
	JobSvc      service.IJobService
	NotifyStore notify.INotifyStore
}

func InitWebServer() *WebContainer {
	wire.Build(
		//database
		database.NewMySQLWithGORM,
		database.NewRedisClient,
		cache.NewEmailCodeCache,
		cache.NewOAuth2StateCache,

		// repo
		repo.NewJobRepo,
		repo.NewJobRecordRepo,
		repo.NewNodeRepo,
		repo.NewUserRepo,
		repo.NewEmailCodeRepo,

		// service
		email.InitEmailService,
		notify.InitMemoryNotifyStore,
		service.NewJobService,
		service.NewJobRecordService,
		service.NewNodeService,
		service.NewEmailCodeService,
		service.NewUserService,
		service.NewDashboardService,

		// api
		api.NewJobApi,
		api.NewJobRecordApi,
		api.NewNodeApi,
		api.NewUserApi,
		api.NewDashboardApi,
		api.NewOAuth2Api,

		// web
		middleware.NewGinMiddlewares,
		router.NewWebRouter,

		wire.Struct(new(WebContainer), "*"),
	)
	return nil
}
