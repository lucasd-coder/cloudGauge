package usageaggregationhistory

import "context"

type (
	UsageAggregationHistoryRepository interface {
		Save(ctx context.Context, usage *UsageAggregationHistory) error
	}
)
