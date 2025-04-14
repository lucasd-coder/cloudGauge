package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation/service"
	"github.com/lucasd-coder/pulseReceiver/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func createServiceWithMocks(t *testing.T) (
	*service.ServiceImpl,
	*mocks.MockValidator,
	*mocks.MockUsageAggregationRepository,
	*mocks.MockUsageAggregationHistoryRepository,
	*mocks.MockPublisher,
	context.Context,
) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockValidator := mocks.NewMockValidator(ctrl)
	mockRepo := mocks.NewMockUsageAggregationRepository(ctrl)
	mockHistoryRepo := mocks.NewMockUsageAggregationHistoryRepository(ctrl)
	mockPublisher := mocks.NewMockPublisher(ctrl)

	svc := &service.ServiceImpl{
		Validate:                          mockValidator,
		Repository:                        mockRepo,
		UsageAggregationHistoryRepository: mockHistoryRepo,
		Publisher:                         mockPublisher,
	}

	ctx := context.Background()

	return svc, mockValidator, mockRepo, mockHistoryRepo, mockPublisher, ctx
}

func TestServiceImpl_Save_Success(t *testing.T) {
	svc, validator, repo, histRepo, _, ctx := createServiceWithMocks(t)

	payload := usageaggregation.Payload{
		Tenant:     "tenant1",
		ProductSKU: "sku1",
		UseUnity:   "unit",
		Amount:     100,
	}

	entities := payload.ToEntities()

	validator.EXPECT().
		ValidateStruct(gomock.Any()).
		Return(nil)

	repo.EXPECT().
		FindToDispatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]usageaggregation.UsageAggregation{}, nil) // corrigido tipo

	histRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	repo.EXPECT().
		Save(gomock.Any(), entities).
		Return(nil)

	err := svc.Save(ctx, payload)
	assert.NoError(t, err)
}

func TestServiceImpl_Save_ValidationError(t *testing.T) {
	svc, validator, _, _, _, ctx := createServiceWithMocks(t)

	payload := usageaggregation.Payload{
		Tenant: "invalid",
	}

	validator.EXPECT().
		ValidateStruct(gomock.AssignableToTypeOf(&usageaggregation.Payload{})).
		Return(errors.New("validation error"))

	err := svc.Save(ctx, payload)
	assert.EqualError(t, err, "validation error")
}

func TestServiceImpl_Save_FindToDispatchError(t *testing.T) {
	svc, validator, repo, _, _, ctx := createServiceWithMocks(t)

	payload := usageaggregation.Payload{
		Tenant:     "tenant",
		ProductSKU: "sku",
		UseUnity:   "unit",
		Amount:     1,
	}

	validator.EXPECT().
		ValidateStruct(gomock.AssignableToTypeOf(&usageaggregation.Payload{})).
		Return(nil)

	repo.EXPECT().
		FindToDispatch(ctx, payload.Tenant, payload.ProductSKU, payload.UseUnity).
		Return(nil, errors.New("find error"))

	err := svc.Save(ctx, payload)
	assert.EqualError(t, err, "find error")
}

func TestServiceImpl_Save_HistorySaveError(t *testing.T) {
	svc, validator, repo, histRepo, _, ctx := createServiceWithMocks(t)

	payload := usageaggregation.Payload{
		Tenant:     "tenant",
		ProductSKU: "sku",
		UseUnity:   "unit",
		Amount:     1,
	}

	validator.EXPECT().
		ValidateStruct(gomock.AssignableToTypeOf(&usageaggregation.Payload{})).
		Return(nil)

	repo.EXPECT().
		FindToDispatch(ctx, payload.Tenant, payload.ProductSKU, payload.UseUnity).
		Return([]usageaggregation.UsageAggregation{}, nil) // <- corrigido

	histRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(errors.New("history save error"))

	err := svc.Save(ctx, payload)
	assert.EqualError(t, err, "history save error")
}

func TestServiceImpl_Save_RepositorySaveError(t *testing.T) {
	svc, validator, repo, histRepo, _, ctx := createServiceWithMocks(t)

	payload := usageaggregation.Payload{
		Tenant:     "tenant",
		ProductSKU: "sku",
		UseUnity:   "unit",
		Amount:     1,
	}

	validator.EXPECT().
		ValidateStruct(gomock.AssignableToTypeOf(&usageaggregation.Payload{})).
		Return(nil)

	repo.EXPECT().
		FindToDispatch(gomock.Any(), gomock.Eq("tenant"), gomock.Eq("sku"), gomock.Eq("unit")).
		Return([]usageaggregation.UsageAggregation{}, nil)

	histRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	repo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(errors.New("save error"))

	err := svc.Save(ctx, payload)
	assert.EqualError(t, err, "save error")
}
