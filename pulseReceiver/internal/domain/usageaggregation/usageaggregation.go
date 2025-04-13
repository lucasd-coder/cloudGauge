package usageaggregation

import (
	"time"

	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

type AggregationStatus string

const (
	StatusPending   AggregationStatus = "pending"
	StatusError     AggregationStatus = "error"
	StatusProcessed AggregationStatus = "processed"
)

type Base struct {
	ID         uint      `gorm:"primaryKey"`
	Tenant     string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	ProductSKU string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	UseUnity   string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	HourWindow time.Time `gorm:"not null;uniqueIndex:idx_aggregation_composite"`

	UsedAmount  float64           `gorm:"not null"`
	Status      AggregationStatus `gorm:"not null;default:'pending'"`
	LastAttempt time.Time
	CreatedAt   time.Time
}

type UsageAggregation struct {
	ID         uint      `gorm:"primaryKey"`
	Tenant     string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	ProductSKU string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	UseUnity   string    `gorm:"not null;uniqueIndex:idx_aggregation_composite"`
	HourWindow time.Time `gorm:"not null;uniqueIndex:idx_aggregation_composite"`

	UsedAmount   float64           `gorm:"not null"`
	Status       AggregationStatus `gorm:"not null;default:'pending'"`
	LastAttempt  time.Time
	CreatedAt    time.Time
	ErrorMessage string
}

type Payload struct {
	Tenant     string  `json:"tenant,omitempty" validate:"required"`
	ProductSKU string  `json:"product_sku,omitempty" validate:"required"`
	UseUnity   string  `json:"use_unity,omitempty" validate:"required"`
	Amount     float64 `json:"used_amount,omitempty" validate:"required"`
}

func (p *Payload) Validate(val shared.Validator) error {
	return val.ValidateStruct(p)
}

func (p *Payload) ToEntities() *UsageAggregation {
	if p == nil {
		return nil
	}

	return &UsageAggregation{
		Tenant:     p.Tenant,
		ProductSKU: p.ProductSKU,
		UseUnity:   p.UseUnity,
		UsedAmount: p.Amount,
		HourWindow: time.Now().Truncate(time.Hour),
		Status:     StatusPending,
	}
}
