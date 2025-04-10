package shared

import (
	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/shared/logger"
)

func NewOptLogger(cfg *config.Config) logger.Option {
	return logger.Option{
		AppName: cfg.Name,
		Level:   cfg.Level,
	}
}

type Options struct {
	URL        string
	GroupId    string
	WorkTask   int
	BatchSize  int
	BatchLimit int
	Poll       int
}
