package notification

import "context"

type SMSService interface {
	SendSMS(ctx context.Context, to, message string) error
}
