package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

type Publisher struct {
	opt     *shared.Options
	publish *kafka.Producer
	cleanup func()
}

func NewPublisher(ctx context.Context, opt *shared.Options) (_ *Publisher, cleanup func(), err error) {
	addCleanup := func(f func()) {
		old := cleanup
		cleanup = func() { old(); f() }
	}

	defer func() {
		if err != nil {
			cleanup()
			cleanup = nil
		}
	}()

	cleanup = func() {}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":         opt.URL,
		"socket.timeout.ms":         10,
		"message.timeout.ms":        10,
		"go.delivery.report.fields": "key,value,headers",
	})
	if err != nil {
		return nil, cleanup, err
	}
	addCleanup(func() {
		p.Close()
	})

	logger.FromContext(ctx).Info("Created", "Producer", p)

	return &Publisher{
		opt:     opt,
		publish: p,
		cleanup: cleanup,
	}, cleanup, nil
}

func (p *Publisher) Publish(ctx context.Context, msg []byte) error {
	deliveryChan := make(chan kafka.Event, 1)

	m := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.opt.TopicName, Partition: kafka.PartitionAny},
		Value:          msg,
	}

	if err := p.publish.Produce(m, deliveryChan); err != nil {
		return fmt.Errorf("error on produce message: %w", err)
	}

	select {
	case e := <-deliveryChan:
		m, ok := e.(*kafka.Message)
		if !ok {
			return fmt.Errorf("unexpected event: %v", e)
		}
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("message delivery failed: %w", m.TopicPartition.Error)
		}
		logger.FromContext(ctx).Infof("Message successfully delivered to %v\n", m.TopicPartition)
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
