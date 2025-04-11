package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/controller"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/middleware"
)

func Start(ctx context.Context, cfg *config.Config) error {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.LoggerMiddleware)

	logger.FromContext(ctx).Infof("Started listening... address[:%s]", cfg.Port)
	controller := controller.NewRouter()

	r.Mount("/", controller)
	r.Mount("/debug", chiMiddleware.Profiler())

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Panic(err)
		return err
	}

	if err := s.Close(); err != nil {
		logger.FromContext(ctx).Error(err.Error())
		return err
	}
	return nil
}
