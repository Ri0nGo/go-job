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

type IOAuth2StateCache interface {
	Set(ctx context.Context, state string, val model.OAuth2State, ttl time.Duration) error
	Get(ctx context.Context, state string) (model.OAuth2State, error)
	MarkUsed(ctx context.Context, state string) error
}
