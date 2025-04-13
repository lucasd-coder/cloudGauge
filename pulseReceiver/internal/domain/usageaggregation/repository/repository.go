package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsageAggregationRepository struct {
	conn *gorm.DB
}

func NewUsageAggregationRepository(conn *gorm.DB) *UsageAggregationRepository {
	return &UsageAggregationRepository{
		conn: conn,
	}
}

func (r *UsageAggregationRepository) Save(ctx context.Context, usage *usageaggregation.UsageAggregation) error {
	tableName := "usage_aggregations"
	expr := fmt.Sprintf(`"%s"."used_amount" + EXCLUDED.used_amount`, tableName)
	return r.conn.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "tenant"}, {Name: "product_sku"}, {Name: "use_unity"}, {Name: "hour_window"}},
			DoUpdates: clause.Assignments(map[string]any{
				"used_amount": gorm.Expr(expr),
				"status":      usageaggregation.StatusPending,
			}),
		}).
		Create(&usage).Error
}

func (r *UsageAggregationRepository) FindToDispatch(
	ctx context.Context,
	tenant, sku, useUnity string,
) ([]usageaggregation.UsageAggregation, error) {
	window := time.Now()
	start := window.Truncate(time.Hour)
	end := start.Add(time.Hour)

	var results []usageaggregation.UsageAggregation

	err := r.conn.WithContext(ctx).
		Where("tenant = ? AND product_sku = ? AND use_unity = ?", tenant, sku, useUnity).
		Where("status IN ?", []usageaggregation.AggregationStatus{
			usageaggregation.StatusPending,
			usageaggregation.StatusError,
		}).
		Where("hour_window >= ? AND hour_window < ?", start, end).
		Find(&results).Error

	return results, err
}
