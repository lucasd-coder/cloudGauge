package usageaggregationhistory

import "time"

type UsageAggregationHistory struct {
	ID         uint      `gorm:"primaryKey"`
	Tenant     string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	ProductSKU string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	UseUnity   string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	HourWindow time.Time `gorm:"not null;uniqueIndex:idx_aggregation_composite"`

	UsedAmount float64 `gorm:"not null"`
	CreatedAt  time.Time
}
