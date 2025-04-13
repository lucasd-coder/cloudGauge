package service

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

var InitializeService = wire.NewSet(
	wire.Bind(new(usageaggregation.UsageAggregationService), new(*ServiceImpl)),
	NewService,
)

type ServiceImpl struct {
	validate                           shared.Validator
	repository                         usageaggregation.UsageAggregationRepository
	usageAggregationHHistoryRepository usageaggregationhistory.UsageAggregationHistoryRepository
}

func NewService(val shared.Validator,
	repo usageaggregation.UsageAggregationRepository,
	historyRepo usageaggregationhistory.UsageAggregationHistoryRepository) *ServiceImpl {
	return &ServiceImpl{
		validate:                           val,
		repository:                         repo,
		usageAggregationHHistoryRepository: historyRepo,
	}
}
