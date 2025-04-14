package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

var (
	producer  *kafka.Producer
	initOnce  sync.Once
	initErr   error
	publisher *Publisher
)

type Publisher struct {
	opt *shared.Options
}

func NewPublisher(opt *shared.Options) (_ *Publisher, err error) {
	initOnce.Do(func() {
		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers":         opt.URL,
			"socket.timeout.ms":         10,
			"message.timeout.ms":        10,
			"go.delivery.report.fields": "key,value,headers",
		})
		if err != nil {
			initErr = fmt.Errorf("failed to kafka.NewProducer: %w", err)
			return
		}

		producer = p
		publisher = &Publisher{
			opt: opt,
		}
	})

	return publisher, initErr
}

func (p *Publisher) Publish(ctx context.Context, msg []byte) error {
	deliveryChan := make(chan kafka.Event, 1)

	m := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.opt.TopicName, Partition: kafka.PartitionAny},
		Value:          msg,
	}

	if err := producer.Produce(m, deliveryChan); err != nil {
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

func CloseProducer() {
	if producer == nil {
		return
	}

	producer.Close()
}
