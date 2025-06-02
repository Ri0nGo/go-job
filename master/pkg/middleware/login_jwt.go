package middleware

import (
	"github.com/Ri0nGo/gokit/slice"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"go-job/master/pkg/config"
	"log/slog"
	"net/http"
	"time"
)

var defaultRefreshJwtTime = time.Minute * 5

type LoginJwtMWBuilder struct {
	key       string
	skipPaths []string // 不需要校验jwt的path
}

func NewLoginJwtMWBuilder(key string) *LoginJwtMWBuilder {
	return &LoginJwtMWBuilder{
		key: key,
	}
}

func (b *LoginJwtMWBuilder) SkipPaths(paths []string) *LoginJwtMWBuilder {
	b.skipPaths = paths
	return b
}

func (b *LoginJwtMWBuilder) isSkipPaths(path string) bool {
	return slice.Contains(b.skipPaths, path)
}

func (b *LoginJwtMWBuilder) Builder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if b.isSkipPaths(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			// 没登录，没有 token, Authorization 这个头部都没有
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 校验token是否有效
		uc := &model.UserClaims{}
		jwtBuilder := auth.NewJwtBuilder(b.key)
		token, err := jwtBuilder.ParseToken(uc, tokenStr)
		if err != nil {
			slog.Error("jwt parse token error", "err", err,
				"token", tokenStr)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 检测token是否需要刷新
		if uc.ExpiresAt.Sub(time.Now()) < defaultRefreshJwtTime {
			jwtToken, err := b.updateJwtToken(jwtBuilder.GetExpireTime(), token, uc)
			ctx.Header("Authorization", jwtToken)
			if err != nil {
				slog.Error("update jwt token error", "err", err, "token", jwtToken)
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		ctx.Set("user", uc)
	}
}

func (b *LoginJwtMWBuilder) updateJwtToken(expireTime time.Duration, token *jwt.Token, uc *model.UserClaims) (string, error) {
	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expireTime))
	newToken, err := token.SignedString([]byte(config.App.Server.Key))
	if err != nil {
		return "", err
	}
	return newToken, nil
}
