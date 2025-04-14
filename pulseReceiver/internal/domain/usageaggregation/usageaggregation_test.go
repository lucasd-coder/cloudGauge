package usageaggregation_test

import (
	"testing"
	"time"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/stretchr/testify/assert"
)

func TestPayload_ToEntities(t *testing.T) {
	now := time.Now().Truncate(time.Hour)

	tests := []struct {
		name     string
		payload  *usageaggregation.Payload
		expected *usageaggregation.UsageAggregation
	}{
		{
			name: "valid payload",
			payload: &usageaggregation.Payload{
				Tenant:     "tenant-1",
				ProductSKU: "sku-123",
				UseUnity:   "unit-xyz",
				Amount:     42.5,
			},
			expected: &usageaggregation.UsageAggregation{
				Tenant:     "tenant-1",
				ProductSKU: "sku-123",
				UseUnity:   "unit-xyz",
				UsedAmount: 42.5,
				Status:     usageaggregation.StatusPending,
			},
		},
		{
			name:     "nil payload",
			payload:  nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payload.ToEntities()

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.Tenant, result.Tenant)
			assert.Equal(t, tt.expected.ProductSKU, result.ProductSKU)
			assert.Equal(t, tt.expected.UseUnity, result.UseUnity)
			assert.Equal(t, tt.expected.UsedAmount, result.UsedAmount)
			assert.Equal(t, tt.expected.Status, result.Status)

			// Comparar HourWindow com tolerância de alguns segundos
			assert.WithinDuration(t, now, result.HourWindow, time.Second*5)
		})
	}
}
