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

type UsageAggregation struct {
	ID         uint      `gorm:"primaryKey" json:"id,omitempty"`
	Tenant     string    `gorm:"not null;uniqueIndex:idx_aggregation_composite" json:"tenant,omitempty"`
	ProductSKU string    `gorm:"not null;uniqueIndex:idx_aggregation_composite" json:"product_sku,omitempty"`
	UseUnity   string    `gorm:"not null;uniqueIndex:idx_aggregation_composite" json:"use_unity,omitempty"`
	HourWindow time.Time `gorm:"not null;uniqueIndex:idx_aggregation_composite" json:"hour_window,omitempty"`

	UsedAmount   float64           `gorm:"not null" json:"used_amount,omitempty"`
	Status       AggregationStatus `gorm:"not null;default:'pending'" json:"-"`
	LastAttempt  time.Time         `json:"last_attempt,omitempty"`
	CreatedAt    time.Time         `json:"created_at,omitempty"`
	ErrorMessage *string           `json:"-"`
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
