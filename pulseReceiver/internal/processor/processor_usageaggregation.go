package processor

import (
	"context"
	"encoding/json"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
)

func (s *Processor) ProcessUsageAggregation(ctx context.Context, msg []byte) error {
	log := logger.FromContext(ctx)

	log.Info("Processing usage aggregation message")
	var payload usageaggregation.Payload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Error("Failed to unmarshal usage aggregation message", "error", err, "raw_message", string(msg))
		return err
	}
	log.Info("Message unmarshaled successfully")
	return s.usageAggregationService.Save(ctx, payload)
}
