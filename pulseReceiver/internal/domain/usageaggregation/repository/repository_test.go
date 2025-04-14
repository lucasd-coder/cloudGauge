package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	tx := db.Begin()

	err = tx.AutoMigrate(&usageaggregation.UsageAggregation{})
	require.NoError(t, err)

	return tx
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	err := db.Rollback().Error
	require.NoError(t, err)
}

func TestUsageAggregationRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := repository.NewUsageAggregationRepository(db)

	usage := &usageaggregation.UsageAggregation{
		Tenant:     "tenant",
		ProductSKU: "sku",
		UseUnity:   "unit",
		HourWindow: time.Now(),
		UsedAmount: 1.0,
		Status:     "pending",
	}

	err := repo.Save(context.Background(), usage)
	require.NoError(t, err)

	var stored usageaggregation.UsageAggregation
	err = db.First(&stored).Error
	require.NoError(t, err)

	require.EqualValues(t, 1.0, stored.UsedAmount)
}

func TestUsageAggregationRepository_FindToDispatch(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := repository.NewUsageAggregationRepository(db)

	usage := &usageaggregation.UsageAggregation{
		Tenant:     "tenant",
		ProductSKU: "sku",
		UseUnity:   "unit",
		HourWindow: time.Now(),
		UsedAmount: 1.0, // float64
		Status:     "pending",
	}

	// Saving usage
	err := repo.Save(context.Background(), usage)
	require.NoError(t, err)

	// Testing FindToDispatch
	result, err := repo.FindToDispatch(context.Background(), "tenant", "sku", "unit")
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.EqualValues(t, "tenant", result[0].Tenant)
	require.EqualValues(t, "sku", result[0].ProductSKU)
	require.EqualValues(t, "unit", result[0].UseUnity)
}

func TestUsageAggregationRepository_FindPending(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := repository.NewUsageAggregationRepository(db)

	usage1 := &usageaggregation.UsageAggregation{
		Tenant:     "tenant1",
		ProductSKU: "sku1",
		UseUnity:   "unit1",
		HourWindow: time.Now(),
		UsedAmount: 1.0,
		Status:     "pending",
	}

	usage2 := &usageaggregation.UsageAggregation{
		Tenant:     "tenant2",
		ProductSKU: "sku2",
		UseUnity:   "unit2",
		HourWindow: time.Now(),
		UsedAmount: 2.0,
		Status:     "error",
	}

	err := repo.Save(context.Background(), usage1)
	require.NoError(t, err)
	err = repo.Save(context.Background(), usage2)
	require.NoError(t, err)

	// Testing FindPending
	result, err := repo.FindPending(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.EqualValues(t, "tenant1", result[0].Tenant)
	require.EqualValues(t, "tenant2", result[1].Tenant)
}

func TestUsageAggregationRepository_UpdateStatusAfterAttempt(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := repository.NewUsageAggregationRepository(db)

	usage := &usageaggregation.UsageAggregation{
		Tenant:     "tenant",
		ProductSKU: "sku",
		UseUnity:   "unit",
		HourWindow: time.Now(),
		UsedAmount: 1.0,
		Status:     "pending",
	}

	// Saving usage
	err := repo.Save(context.Background(), usage)
	require.NoError(t, err)

	// Testing UpdateStatusAfterAttempt (success)
	err = repo.UpdateStatusAfterAttempt(context.Background(), usage, true, nil)
	require.NoError(t, err)

	var updated usageaggregation.UsageAggregation
	err = db.First(&updated).Error
	require.NoError(t, err)

	require.EqualValues(t, usageaggregation.StatusProcessed, updated.Status)
	require.Nil(t, updated.ErrorMessage) // Verifica se é nil, não uma string vazia

	// Testing UpdateStatusAfterAttempt (failure)
	errMessage := "some error"
	err = repo.UpdateStatusAfterAttempt(context.Background(), usage, false, &errMessage)
	require.NoError(t, err)

	err = db.First(&updated).Error
	require.NoError(t, err)

	require.EqualValues(t, usageaggregation.StatusError, updated.Status)
	require.EqualValues(t, errMessage, *updated.ErrorMessage) // Verifica o valor da mensagem de erro
}
