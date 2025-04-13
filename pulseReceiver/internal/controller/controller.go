package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	appError "github.com/lucasd-coder/pulseReceiver/internal/shared/errors"
)

type controller struct{}

func NewRouter(
	usageAggregationController *UsageAggregationController) *chi.Mux {
	r := chi.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	r.Group(func(r chi.Router) {
		r.Route("/usage", func(r chi.Router) {
			r.Post("/", usageAggregationController.Save)
		})
	})

	return r
}

func (c *controller) SendError(ctx context.Context, w http.ResponseWriter, err error) {
	errResp := appError.BuildError(err)

	c.Response(ctx, w, errResp, errResp.StatusCode)
}

func (c *controller) Response(ctx context.Context, w http.ResponseWriter, body interface{}, statusCode int) {
	log := logger.FromContext(ctx)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	content, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		log.Error("err during json.Marchal", err)
	}

	if _, err := w.Write(content); err != nil {
		log.Error("err during http.ResponseWriter", err)
	}
}
