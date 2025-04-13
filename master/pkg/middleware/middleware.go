package middleware

import (
	"github.com/gin-gonic/gin"
	"go-job/master/pkg/config"
)

func NewGinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors(),
		NewLoginJwtMWBuilder(config.App.Server.Key).SkipPaths([]string{
			"/api/go-job/users/login",
		}).Builder(),
	}
}
