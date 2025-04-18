package email

import "context"

type EmailService interface {
	Send(ctx context.Context, email string) error
}
