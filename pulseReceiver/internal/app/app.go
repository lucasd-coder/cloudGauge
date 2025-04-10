package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/controller"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
	"github.com/lucasd-coder/pulseReceiver/internal/shared/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/shared/middleware"
)

func Run(cfg *config.Config) {
	optlogger := shared.NewOptLogger(cfg)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger.NewLogger(optlogger)
	logDefault := logger.GetLog()
	slog.SetDefault(logDefault)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.LoggerMiddleware)

	logDefault.Info(fmt.Sprintf("Started listening... address[:%s]", cfg.Port))

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
		return
	}

	if err := s.Close(); err != nil {
		logDefault.Error(err.Error())
		return
	}
}
