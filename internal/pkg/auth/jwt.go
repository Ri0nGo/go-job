package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-job/internal/model"
	"go-job/internal/pkg/consts"
	"time"
)

const (
	defaultExpireTime = time.Hour * 2
)

var (
	ErrorTokenInvalid = errors.New("token is invalid")
	ErrorTokenExpired = errors.New("token is expired")
	ErrorTokenIsNil   = errors.New("token is nil")
)

type JwtBuilder struct {
	key []byte
}

func NewJwtBuilder(key string) *JwtBuilder {
	builder := &JwtBuilder{
		key: []byte(key),
	}
	return builder
}

// GenerateToken 生成jwt token
// 注意：过期时间是由调用者设置的
func (builder *JwtBuilder) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(builder.key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// GenerateUserToken 生成用户登录的jwt token
func (builder *JwtBuilder) GenerateUserToken(user model.DomainUser) (string, error) {
	uc := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.DefaultLoginJwtExpireTime)),
		},
		Uid: user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(builder.key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseToken 解析jwt token
func (builder *JwtBuilder) ParseToken(claims jwt.Claims, tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return builder.key, nil
	})
	if err != nil {
		return nil, ErrorTokenInvalid
	}
	if token == nil || !token.Valid {
		return nil, ErrorTokenExpired
	}
	return token, nil
}
