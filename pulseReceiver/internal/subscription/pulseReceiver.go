package subscription

import (
	"context"

	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/inject"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/kafka"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

func PulseReceiverEvent(ctx context.Context) error {
	cfg := config.GetConfig()

	topicName := "aggregated-pulses"
	opt := shared.NewOptions(cfg, topicName)

	p := inject.InitializeProcessor()

	start, cleanup, err := kafka.NewSubscription(ctx, opt, p.ProcessUsageAggregation)
	if err != nil {
		return err
	}
	defer cleanup()
	if err := start.Start(ctx); err != nil {
		return err
	}
	return nil
}
