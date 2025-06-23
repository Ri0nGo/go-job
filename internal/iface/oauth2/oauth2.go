package oauth2

import (
	"context"
	"go-job/internal/model"
	"time"
)

type IOAuth2Service interface {
	GetAuthUrl(ctx context.Context, state string) string
	GetAuthIdentity(ctx context.Context, code string) (model.AuthIdentity, error)
}

type IOAuth2Cache interface {
	SetState(ctx context.Context, state string, oauth2State model.OAuth2State) error
	SetAuth(ctx context.Context, key string, val map[string]string, ttl time.Duration) error
	GetState(ctx context.Context, state string) (model.OAuth2State, error)
	GetAuth(ctx context.Context, key string) (map[string]string, error)
	MarkUsed(ctx context.Context, state string, flag model.OAuth2Flag) error
}
