package usageaggregationhistory

import "time"

type UsageAggregationHistory struct {
	ID         uint      `gorm:"primaryKey"`
	Tenant     string    `gorm:"not null"`
	ProductSKU string    `gorm:"not null"`
	UseUnity   string    `gorm:"not null"`
	HourWindow time.Time `gorm:"not null"`

	UsedAmount float64 `gorm:"not null"`
	CreatedAt  time.Time
}
