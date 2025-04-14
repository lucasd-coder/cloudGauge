package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
)

func (s *ServiceImpl) Save(ctx context.Context, payload usageaggregation.Payload) error {
	log := logger.FromContext(ctx)
	log.Info("Received usage aggregation event", "payload", payload)

	if err := payload.Validate(s.Validate); err != nil {
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

	if err := s.UsageAggregationHistoryRepository.Save(ctx, &history); err != nil {
		log.Error("Failed to save usage aggregation history", "error", err)
		return err
	}

	log.Info("Saving usage aggregation entities", "entities", aggregation)
	if err := s.Repository.Save(ctx, payload.ToEntities()); err != nil {
		log.Error("Failed to save usage aggregation entities", "error", err)
		return err
	}

	log.Info("Usage aggregation event successfully saved")
	return nil
}

func (s *ServiceImpl) findToDispatch(ctx context.Context, payload usageaggregation.Payload) (*usageaggregation.UsageAggregation, error) {
	usages, err := s.Repository.FindToDispatch(ctx, payload.Tenant, payload.ProductSKU, payload.UseUnity)
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

func (s *ServiceImpl) DispatchPending(ctx context.Context) {
	log := logger.FromContext(ctx)

	now := time.Now()
	window := now.Truncate(time.Hour)

	log.Info("Starting hourly dispatch job", "job_hour_window", window)

	records, err := s.Repository.FindPending(ctx)
	if err != nil {
		log.Error("Failed to query usage aggregations for dispatch", "error", err)
		return
	}

	if len(records) == 0 {
		log.Info("No pending usage aggregations to dispatch", "job_hour_window", window)
		return
	}

	for _, record := range records {
		jsonData, err := json.Marshal(record)
		if err != nil {
			log.Error("Failed to marshal usage aggregation",
				"error", err,
				"tenant", record.Tenant,
				"sku", record.ProductSKU,
				"unity", record.UseUnity,
				"record_hour_window", record.HourWindow)
			continue
		}

		err = s.Publisher.Publish(ctx, jsonData)
		if err != nil {
			log.Error("Failed to publish usage aggregation",
				"error", err,
				"tenant", record.Tenant,
				"sku", record.ProductSKU,
				"unity", record.UseUnity,
				"record_hour_window", record.HourWindow)

			msg := err.Error()
			if err := s.Repository.UpdateStatusAfterAttempt(ctx, &record, false, &msg); err != nil {
				log.Error("Failed to update usage aggregation status after publish failure",
					"error", err,
					"tenant", record.Tenant,
					"sku", record.ProductSKU,
					"unity", record.UseUnity,
					"record_hour_window", record.HourWindow)
			}
			continue
		}

		if err := s.Repository.UpdateStatusAfterAttempt(ctx, &record, true, nil); err != nil {
			log.Error("Failed to update usage aggregation status after success",
				"error", err,
				"tenant", record.Tenant,
				"sku", record.ProductSKU,
				"unity", record.UseUnity,
				"record_hour_window", record.HourWindow)
		}

		log.Info("Successfully dispatched usage aggregation",
			"tenant", record.Tenant,
			"sku", record.ProductSKU,
			"unity", record.UseUnity,
			"record_hour_window", record.HourWindow)
	}
}
