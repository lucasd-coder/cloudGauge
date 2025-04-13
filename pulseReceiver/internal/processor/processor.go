package processor

import (
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
)

type Processor struct {
	usageAggregationService usageaggregation.UsageAggregationService
}

func NewProcessor(usageAggregationService usageaggregation.UsageAggregationService) *Processor {
	return &Processor{usageAggregationService: usageAggregationService}
}
