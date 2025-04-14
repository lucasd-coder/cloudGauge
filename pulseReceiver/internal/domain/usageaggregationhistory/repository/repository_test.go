package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregationhistory/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Configura um banco de dados SQLite em memória
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Realiza a migração para garantir que as tabelas sejam criadas
	err = db.AutoMigrate(&usageaggregationhistory.UsageAggregationHistory{})
	require.NoError(t, err)

	return db
}

func TestUsageAggregationHistoryRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUsageAggregationHistoryRepository(db)

	usage := &usageaggregationhistory.UsageAggregationHistory{
		Tenant:     "tenant",
		ProductSKU: "sku",
		UseUnity:   "unit",
		HourWindow: time.Now(),
		UsedAmount: 1.0, // float64
	}

	err := repo.Save(context.Background(), usage)
	require.NoError(t, err)

	var stored usageaggregationhistory.UsageAggregationHistory
	err = db.First(&stored, "tenant = ? AND product_sku = ?", "tenant", "sku").Error
	require.NoError(t, err)

	require.Equal(t, "tenant", stored.Tenant)
	require.Equal(t, "sku", stored.ProductSKU)
	require.Equal(t, "unit", stored.UseUnity)
	require.Equal(t, usage.UsedAmount, stored.UsedAmount)
}
