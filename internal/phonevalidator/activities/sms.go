package activities

import (
	"context"
	"errors"
	"fmt"

	"go.temporal.io/sdk/activity"
)

//go:generate options-gen -out-filename=sms_options.gen.go -from-struct=SMSOptions
type SMSOptions struct {
	defaultGateway string `option:"mandatory" validate:"required"`
}

type SMSActivity struct {
	SMSOptions
	gateways map[string]SMSGateway
}

func NewSMSActivity(opts SMSOptions) (*SMSActivity, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate sms activity: %v", err)
	}

	return &SMSActivity{
		SMSOptions: opts,
		gateways:   make(map[string]SMSGateway),
	}, nil
}

// RegisterGateway registers a new SMS gateway with the given ID.
//
// smsGatewayID: the ID of the SMS gateway
// gateway: the SMS gateway to register
// Method is not thread safe.
func (a *SMSActivity) RegisterGateway(smsGatewayID string, gateway SMSGateway) {
	a.gateways[smsGatewayID] = gateway
}

func (a *SMSActivity) SendCodeViaSMS(ctx context.Context, phone, msg string) error {
	l := activity.GetLogger(ctx)

	gw, ok := a.gateways[a.defaultGateway]
	if !ok {
		return errors.New("no default sms gateway")
	}

	if err := gw.SendSMS(ctx, phone, msg); err != nil {
		l.Error("Send sms", "Error", err)

		return fmt.Errorf("send sms: %v", err)
	}

	l.Info("Success send via sms")
	return nil
}
