package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"log/slog"
	"sync"
	"time"
)

const (
	defaultExpireTime       = time.Hour * 2
	defaultRefreshTimeLimit = time.Minute * 5
	defaultJwtExpireTime    = time.Hour * 2
	defaultUserId           = -1
)

type JwtToken struct {
	builder     *auth.JwtBuilder // TODO 这里可以改成接口，不同的校验都去实现生成和校验接口，这样就热插拔的方式了
	token       string
	ExpireTime  time.Time     `json:"expire_time"`  // 过期时间
	RefreshTime time.Duration `json:"refresh_time"` // 过期前多久刷新
	mux         sync.RWMutex
}

var jt = &JwtToken{}

func InitJwtToken(key string) {
	jt.builder = auth.NewJwtBuilder(key)
	jt.RefreshTime = defaultRefreshTimeLimit
}

func RefreshToken() error {
	jt.mux.Lock()
	defer jt.mux.Unlock()
	// 初始化token， 快过期刷新token
	if jt.token == "" || jt.ExpireTime.Sub(time.Now()) < jt.RefreshTime {
		jt.ExpireTime = time.Now().Add(defaultExpireTime)
		uc := model.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(jt.ExpireTime),
			},
			Uid: defaultUserId,
		}
		token, err := jt.builder.GenerateToken(uc)
		if err != nil {
			return err
		}
		jt.token = token
		slog.Info("refresh token success", "token", jt.token,
			"expire_time", jt.ExpireTime, "refresh_time", jt.RefreshTime)
	}
	return nil
}

func GetJwtToken() string {
	jt.mux.RLock()
	defer jt.mux.RUnlock()
	return jt.token
}
