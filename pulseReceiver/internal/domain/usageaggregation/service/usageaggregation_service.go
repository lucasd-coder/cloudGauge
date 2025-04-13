package service

import (
	"context"
	"time"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
)

func (s *ServiceImpl) Save(ctx context.Context, payload usageaggregation.Payload) error {
	log := logger.FromContext(ctx)
	log.Info("Received usage aggregation event", "payload", payload)

	if err := payload.Validate(s.validate); err != nil {
		log.Error("Payload validation failed", "error", err, "payload", payload)
		return err
	}

	aggregation, err := s.findToDispatch(ctx, payload)
	if err != nil {
		return err
	}

	history := usageaggregationhistory.UsageAggregationHistory{
		Tenant:     payload.Tenant,
		ProductSKU: payload.ProductSKU,
		UseUnity:   payload.UseUnity,
		UsedAmount: payload.Amount,
		HourWindow: time.Now().Truncate(time.Hour),
	}

	log.Info("Saving usage aggregation history", "history_base", history)

	if err := s.usageAggregationHHistoryRepository.Save(ctx, &history); err != nil {
		log.Error("Failed to save usage aggregation history", "error", err)
		return err
	}

	log.Info("Saving usage aggregation entities", "entities", aggregation)
	if err := s.repository.Save(ctx, payload.ToEntities()); err != nil {
		log.Error("Failed to save usage aggregation entities", "error", err)
		return err
	}

	log.Info("Usage aggregation event successfully saved")
	return nil
}

func (s *ServiceImpl) findToDispatch(ctx context.Context, payload usageaggregation.Payload) (*usageaggregation.UsageAggregation, error) {
	usages, err := s.repository.FindToDispatch(ctx, payload.Tenant, payload.ProductSKU, payload.UseUnity)
	if err != nil {
		return nil, err
	}
	entities := payload.ToEntities()
	if len(usages) == 0 {
		return entities, nil
	}

	for _, usage := range usages {
		entities.ID = usage.ID
		entities.UsedAmount += usage.UsedAmount
	}

	return entities, nil
}
