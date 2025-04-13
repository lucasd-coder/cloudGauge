package repository

import (
	"context"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	"gorm.io/gorm"
)

type UsageAggregationHistoryRepository struct {
	conn *gorm.DB
}

func NewUsageAggregationHistoryRepository(conn *gorm.DB) *UsageAggregationHistoryRepository {
	return &UsageAggregationHistoryRepository{
		conn: conn,
	}
}

func (r *UsageAggregationHistoryRepository) Save(ctx context.Context, usage *usageaggregationhistory.UsageAggregationHistory) error {
	return r.conn.Create(&usage).Error
}
