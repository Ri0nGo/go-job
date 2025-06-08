package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-job/internal/pkg/ratelimit"
	"go-job/master/pkg/config"
	"go-job/master/pkg/middleware/ratelimit/iplimiter"
	"time"
)

const (
	defaultLimitInterval = time.Second * 1
	defaultLimitRate     = 50
)

func NewGinMiddlewares(cmd redis.Cmdable) []gin.HandlerFunc {
	redisLimiter := ratelimit.NewRedisSlidingWindowLimiter(cmd, defaultLimitInterval, defaultLimitRate)
	return []gin.HandlerFunc{
		cors(),
		NewLoginJwtMWBuilder(config.App.Server.Key).SkipPaths([]string{
			"/api/go-job/users/login",
			"/api/go-job/oauth2/github/authurl",
			"/api/go-job/oauth2/github/callback",
		}).Builder(),
		iplimiter.NewIpLimiter(redisLimiter).Builder(),
	}
}
