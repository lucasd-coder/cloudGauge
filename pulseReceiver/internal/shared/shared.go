package shared

import (
	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
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
	TopicName  string
}

func NewOptions(cfg *config.Config, topicName string) *Options {
	return &Options{
		URL:        cfg.Kafka.URL,
		GroupId:    cfg.GroupId,
		WorkTask:   cfg.WorkTask,
		BatchSize:  cfg.BatchSize,
		BatchLimit: cfg.BatchLimit,
		Poll:       cfg.Poll,
		TopicName:  topicName,
	}
}
