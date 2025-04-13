package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/postgres"
	"github.com/lucasd-coder/pulseReceiver/internal/server"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
	"github.com/lucasd-coder/pulseReceiver/internal/subscription"
)

func Run(cfg *config.Config) {
	optlogger := shared.NewOptLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger.NewLogger(optlogger)
	logDefault := logger.GetLog()
	slog.SetDefault(logDefault)
	var wg sync.WaitGroup

	// Postgres config
	postgres.StartDB(ctx, cfg)

	// Postgres closes
	defer postgres.CloseConn()

	// starting the server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Start(ctx, cfg); err != nil {
			log.Fatal(err)
		}
	}()

	// starting the subscriptions
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := subscription.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// channel to lister for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// wait for a signals
	<-sigChan
	logDefault.Info("Received shutdown signal")

	// cancel the context to signal all goroutine to stop
	cancel()

	// wait for all goroutines to complete
	wg.Wait()

	logDefault.Info("shutdown complete")
}
