package iplimiter

import (
	"github.com/gin-gonic/gin"
	"go-job/internal/pkg/ratelimit"
	"net/http"
)

/*
基于redis实现的Ip的滑动窗口限流
*/

type IpLimitBuilder struct {
	prefix  string
	limiter ratelimit.Limiter
}

func NewIpLimiter(limiter ratelimit.Limiter) *IpLimitBuilder {
	return &IpLimitBuilder{
		prefix:  "ip-limiter",
		limiter: limiter,
	}
}

func (b *IpLimitBuilder) SetPrefix(prefix string) *IpLimitBuilder {
	b.prefix = prefix
	return b
}

func (b *IpLimitBuilder) Builder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := b.genKey(ctx)
		isLimit, err := b.limiter.Limit(ctx, key)
		if err != nil { // 表示redis有问题了
			// 根据业务的情况，是用户优先还是业务稳定性邮箱
			// 用户优先则继续让用户通过
			// 这里采取保守策略，不允许用户再访问了
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if isLimit {
			// todo 这里后续可以记录一下是哪些client的ip
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *IpLimitBuilder) genKey(ctx *gin.Context) string {
	return b.prefix + ctx.ClientIP()
}
