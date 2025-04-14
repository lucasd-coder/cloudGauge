package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealthCheck(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestController_Response(t *testing.T) {
	ctrl := &controller{}
	w := httptest.NewRecorder()
	ctx := context.Background()

	body := map[string]string{"message": "success"}

	ctrl.Response(ctx, w, body, http.StatusCreated)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{
		"message": "success"
	}`, w.Body.String())
}

func TestController_SendError(t *testing.T) {
	ctrl := &controller{}
	w := httptest.NewRecorder()
	ctx := context.Background()

	testErr := errors.New("something went wrong")
	ctrl.SendError(ctx, w, testErr)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "something went wrong")
}

type mockUsageService struct {
	mock.Mock
}

func (m *mockUsageService) Save(ctx context.Context, payload usageaggregation.Payload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func TestUsageAggregation_Save_InvalidJSON(t *testing.T) {
	mockService := new(mockUsageService)
	ctrl := NewUsageAggregationController(mockService)

	router := chi.NewRouter()
	router.Post("/usage", ctrl.Save)

	body := bytes.NewBufferString("invalid-json")
	req := httptest.NewRequest(http.MethodPost, "/usage", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error when doing decoder payload")
}

func TestUsageAggregation_Save_Success(t *testing.T) {
	mockService := new(mockUsageService)

	expected := usageaggregation.Payload{
		Tenant: "123",
		Amount: 50,
	}

	mockService.
		On("Save", mock.Anything, expected).
		Return(nil)

	ctrl := NewUsageAggregationController(mockService)

	router := chi.NewRouter()
	router.Post("/usage", ctrl.Save)

	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPost, "/usage", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Body.String()) // body vazio como no controller
}
