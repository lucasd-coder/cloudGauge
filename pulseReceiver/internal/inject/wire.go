//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/pulseReceiver/internal/controller"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	usageaggregationrepository "github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation/repository"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation/service"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	usageaggregationhistoryrepository "github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory/repository"
	"github.com/lucasd-coder/pulseReceiver/internal/processor"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/postgres"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/validator"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

var initializeValidator = wire.NewSet(
	wire.Struct(new(validator.Validation)),
	wire.Bind(new(shared.Validator), new(*validator.Validation)),
)

var initializeUsageAggregationRepository = wire.NewSet(
	wire.Bind(new(usageaggregation.UsageAggregationRepository), new(*usageaggregationrepository.UsageAggregationRepository)),
	usageaggregationrepository.NewUsageAggregationRepository,
)

var initializeUsageAggregationHistoryRepository = wire.NewSet(
	wire.Bind(new(usageaggregationhistory.UsageAggregationHistoryRepository), new(*usageaggregationhistoryrepository.UsageAggregationHistoryRepository)),
	usageaggregationhistoryrepository.NewUsageAggregationHistoryRepository,
)

func InitializeProcessor() *processor.Processor {
	wire.Build(initializeUsageAggregationRepository, initializeUsageAggregationHistoryRepository, initializeValidator, service.InitializeService, postgres.GetConn, processor.NewProcessor)
	return &processor.Processor{}
}

func InitializeUsageAggregationController() *controller.UsageAggregationController {
	wire.Build(initializeUsageAggregationRepository, initializeUsageAggregationHistoryRepository, initializeValidator, service.InitializeService, postgres.GetConn, controller.NewUsageAggregationController)
	return &controller.UsageAggregationController{}
}
