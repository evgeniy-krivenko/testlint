package phonevalidator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"golang.org/x/sync/errgroup"

	"testlint/internal/phonevalidator/activities"
	"testlint/internal/phonevalidator/gateway/smsru"
)

//go:generate options-gen -out-filename=service_opts.gen.go -from-struct=Options
type Options struct {
	temporalClient client.Client `option:"mandatory" validate:"required"`
	callProvider   string        `option:"mandatory" validate:"required"`
	smsProvider    string        `option:"mandatory" validate:"required"`

	poolingTimeout        time.Duration `default:"100ms"`
	workflowExecTimeout   time.Duration `default:"600s"`
	defaultRequestTimeout time.Duration `default:"5s"`
}

type Service struct {
	Options
	callActivity *activities.CallActivity
	smsActivity  *activities.SMSActivity
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate phone validator service: %v", err)
	}

	gw := new(smsru.Gateway)

	callActivity, err := activities.NewCallActivity(activities.NewCallOptions(opts.callProvider))
	if err != nil {
		return nil, fmt.Errorf("create call activity: %v", err)
	}

	smsActivity, err := activities.NewSMSActivity(activities.NewSMSOptions(opts.smsProvider))
	if err != nil {
		return nil, fmt.Errorf("create sms activity: %v", err)
	}

	callActivity.RegisterGateway(opts.callProvider, gw)
	smsActivity.RegisterGateway(opts.smsProvider, gw)

	if err != nil {
		return nil, fmt.Errorf("create validation workflow: %v", err)
	}

	return &Service{
		Options:      opts,
		callActivity: callActivity,
		smsActivity:  smsActivity,
	}, nil
}

// Run runs the service with the provided context.
func (s *Service) Run(ctx context.Context) error {
	w := worker.New(s.temporalClient, "VALIDATION_TASK_QUEUE", worker.Options{})

	w.RegisterActivity(s.callActivity.SendCodeViaCall)
	w.RegisterActivity(s.smsActivity.SendCodeViaSMS)
	w.RegisterActivity(activities.GenerateCode)

	eg, ctx := errgroup.WithContext(ctx)
	doneCh := make(chan interface{})

	eg.Go(func() error {
		<-ctx.Done()
		close(doneCh)
		return nil
	})

	eg.Go(func() error { return w.Run(worker.InterruptCh()) })

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop phone validator: %v", err)
	}

	return nil
}
