package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-job/internal/model"
	"go-job/master/pkg/config"
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

type Option func(builder *JwtBuilder)

type JwtBuilder struct {
	key        []byte
	expireTime time.Duration // 过期时间，单位秒
}

func NewJwtBuilder(key string, opts ...Option) *JwtBuilder {
	builder := &JwtBuilder{
		expireTime: defaultExpireTime,
		key:        []byte(key),
	}
	for _, opt := range opts {
		opt(builder)
	}
	return builder
}

func WithExpireTime(expireTime time.Duration) Option {
	return func(builder *JwtBuilder) {
		if expireTime > 0 {
			builder.expireTime = expireTime
		}
	}
}

// GenerateToken 生成jwt token
func (builder *JwtBuilder) GenerateToken(user model.DomainUser) (string, error) {
	uc := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(builder.expireTime)),
		},
		Uid: user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString([]byte(builder.key))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseToken 解析jwt token
func (builder *JwtBuilder) ParseToken(uc *model.UserClaims, tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.App.Server.Key), nil
	})
	if err != nil {
		return nil, ErrorTokenInvalid
	}
	if token == nil || !token.Valid {
		return nil, ErrorTokenExpired
	}
	return token, nil
}

func (builder *JwtBuilder) GetExpireTime() time.Duration {
	return builder.expireTime
}
