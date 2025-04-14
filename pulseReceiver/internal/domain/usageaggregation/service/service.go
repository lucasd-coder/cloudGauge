package service

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/kafka"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

var InitializeService = wire.NewSet(
	wire.Bind(new(usageaggregation.UsageAggregationService), new(*ServiceImpl)),
	NewService,
)

type ServiceImpl struct {
	Validate                          shared.Validator
	Repository                        usageaggregation.UsageAggregationRepository
	UsageAggregationHistoryRepository usageaggregationhistory.UsageAggregationHistoryRepository
	Publisher                         shared.Publisher
}

func NewService(val shared.Validator,
	repo usageaggregation.UsageAggregationRepository,
	historyRepo usageaggregationhistory.UsageAggregationHistoryRepository) *ServiceImpl {
	cfg := config.GetConfig()

	topicName := "aggregated-to-process"
	opt := shared.NewOptions(cfg, topicName)

	publisher, err := kafka.NewPublisher(opt)
	if err != nil {
		slog.Error("Error on kafka.NewPublisher", "error", err)
		return nil
	}
	return &ServiceImpl{
		Validate:                          val,
		Repository:                        repo,
		UsageAggregationHistoryRepository: historyRepo,
		Publisher:                         publisher,
	}
}
