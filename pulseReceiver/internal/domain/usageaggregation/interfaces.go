package usageaggregation

import (
	"context"
)

type (
	UsageAggregationRepository interface {
		Save(ctx context.Context, usage *UsageAggregation) error
		FindToDispatch(
			ctx context.Context,
			tenant, sku, useUnity string,
		) ([]UsageAggregation, error)
		FindPending(ctx context.Context) ([]UsageAggregation, error)
		UpdateStatusAfterAttempt(
			ctx context.Context,
			usage *UsageAggregation,
			success bool,
			errorMessage *string,
		) error
	}

	UsageAggregationService interface {
		Save(ctx context.Context, payload Payload) error
	}
)
