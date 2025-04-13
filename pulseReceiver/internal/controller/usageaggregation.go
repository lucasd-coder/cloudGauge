package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lucasd-coder/pulseReceiver/internal/domain/usageaggregation"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
)

type UsageAggregationController struct {
	controller
	usageAggregationService usageaggregation.UsageAggregationService
}

func NewUsageAggregationController(
	usageAggregationService usageaggregation.UsageAggregationService) *UsageAggregationController {
	return &UsageAggregationController{
		usageAggregationService: usageAggregationService,
	}
}

func (c *UsageAggregationController) Save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)

	pld := &usageaggregation.Payload{}

	if err := json.NewDecoder(r.Body).Decode(pld); err != nil {
		msg := fmt.Errorf("error when doing decoder payload: %w", err)
		log.Error(msg.Error())
		c.SendError(ctx, w, msg)
		return
	}

	if err := c.usageAggregationService.Save(ctx, *pld); err != nil {
		c.SendError(ctx, w, err)
		return
	}

	c.Response(ctx, w, nil, http.StatusOK)
}
