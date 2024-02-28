package activities

import (
	"context"
	"errors"
	"fmt"

	"go.temporal.io/sdk/activity"
)

//go:generate options-gen -out-filename=call_options.gen.go -from-struct=CallOptions
type CallOptions struct {
	defaultGateway string `option:"mandatory" validate:"required"`
}

type CallActivity struct {
	CallOptions
	gateways map[string]CallGateway
}

func NewCallActivity(opts CallOptions) (*CallActivity, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate call activity: %v", err)
	}

	return &CallActivity{
		CallOptions: opts,
		gateways:    make(map[string]CallGateway),
	}, nil
}

func (a *CallActivity) RegisterGateway(callGatewayID string, gateway CallGateway) {
	a.gateways[callGatewayID] = gateway
}

func (a *CallActivity) SendCodeViaCall(ctx context.Context, phone string) (string, error) {
	l := activity.GetLogger(ctx)

	gw, ok := a.gateways[a.defaultGateway]
	if !ok {
		return "", errors.New("no default call gateway")
	}

	code, err := gw.Call(ctx, phone)
	if err != nil {
		l.Error("Call gateway", "Error", err)

		return "", fmt.Errorf("call gateway: %v", err)
	}

	l.Info("Success send via call")
	return code, nil
}
