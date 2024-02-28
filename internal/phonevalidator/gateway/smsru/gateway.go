package smsru

import (
	"context"

	"testlint/internal/phonevalidator/activities"
)

var _ activities.CallGateway = (*Gateway)(nil)
var _ activities.SMSGateway = (*Gateway)(nil)

type Gateway struct{}

func (g *Gateway) SendSMS(ctx context.Context, phone, msg string) error {
	return nil
}

func (g *Gateway) Call(ctx context.Context, phone string) (code string, err error) {
	return "1234", nil
}
