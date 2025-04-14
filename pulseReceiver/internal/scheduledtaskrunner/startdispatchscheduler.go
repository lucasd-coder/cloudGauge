package scheduledtaskrunner

import (
	"context"

	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/inject"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/jobscheduled"
)

func StartDispatchScheduler(ctx context.Context) error {
	cfg := config.GetConfig()

	dispatch := inject.InitializeUsageAggregationService()

	start, cleanup, err := jobscheduled.NewJobScheduler(ctx, dispatch.DispatchPending, cfg.AggregateCurrentHour)
	if err != nil {
		return err
	}
	defer cleanup()
	if err := start.Start(ctx); err != nil {
		return err
	}
	return nil
}
