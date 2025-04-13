package auth

import (
	"go-job/internal/model"
	"go-job/internal/pkg/auth"
	"time"
)

const (
	defaultExpireTime       = time.Hour * 2
	defaultRefreshTimeLimit = time.Minute * 5
	defaultUserId           = -1
)

//var uc = model.UserClaims{
//	RegisteredClaims: jwt.RegisteredClaims{
//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultExpireTime)),
//	},
//	Uid: defaultUserId,
//}

var uc = model.DomainUser{
	Id: defaultUserId,
}

type JwtToken struct {
	builder     *auth.JwtBuilder // TODO 这里可以改成接口，不同的校验都去实现生成和校验接口，这样就热插拔的方式了
	token       string
	ExpireTime  time.Time     `json:"expire_time"`  // 过期时间
	RefreshTime time.Duration `json:"refresh_time"` // 过期前多久刷新
}

var jt = &JwtToken{}

func InitJwtToken(key string) {
	jt.builder = auth.NewJwtBuilder(key)
}

func RefreshToken() error {
	// 初始化token， 快过期刷新token
	if jt.token == "" || jt.ExpireTime.Sub(time.Now()) < jt.RefreshTime {
		token, err := jt.builder.GenerateToken(uc)
		if err != nil {
			return err
		}
		jt.token = token

	}
	return nil
}

func GetJwtToken() string {
	return jt.token
}
