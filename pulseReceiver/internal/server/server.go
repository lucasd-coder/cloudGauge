package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/controller"
	"github.com/lucasd-coder/pulseReceiver/internal/inject"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/middleware"
)

func Start(ctx context.Context, cfg *config.Config) error {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.LoggerMiddleware)

	log := logger.FromContext(ctx)

	usageAggregationController := inject.InitializeUsageAggregationController()
	controller := controller.NewRouter(usageAggregationController)

	r.Mount("/", controller)
	r.Mount("/debug", chiMiddleware.Profiler())

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	wait := time.Second * 15

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Errorf("Server has stopped due to %v\n", err)
		}
	}()
	log.Infof("Started listening... address[:%s]", cfg.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// wait signal
	<-c
	ctx, cancel := context.WithTimeout(ctx, wait)
	defer cancel()
	return s.Shutdown(ctx)
}
