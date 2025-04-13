package migrations

import (
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(usageaggregation.UsageAggregation{},
		usageaggregationhistory.UsageAggregationHistory{})
}
