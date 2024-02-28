package activities

import "context"

type SMSGateway interface {
	SendSMS(ctx context.Context, phone, msg string) error
}

type CallGateway interface {
	Call(ctx context.Context, phone string) (code string, err error)
}
