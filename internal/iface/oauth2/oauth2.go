package oauth2

import (
	"context"
	"go-job/internal/model"
)

type IOAuth2Service interface {
	GetAuthUrl(ctx context.Context, state string) string
	GetAuthIdentity(ctx context.Context, code string) (model.AuthIdentity, error)
}
