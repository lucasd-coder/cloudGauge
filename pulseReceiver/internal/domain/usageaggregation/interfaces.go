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
	}

	UsageAggregationService interface {
		Save(ctx context.Context, payload Payload) error
	}
)
