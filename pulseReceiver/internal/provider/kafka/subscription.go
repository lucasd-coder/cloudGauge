package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/shared"
)

type Subscription struct {
	handler  func(ctx context.Context, msg []byte) error
	msg      chan *kafka.Message
	opt      *shared.Options
	consumer *kafka.Consumer
	cleanup  func()
}

func NewSubscription(
	ctx context.Context,
	opt *shared.Options,
	h func(ctx context.Context, msg []byte) error) (_ *Subscription, cleanup func(), err error) {
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

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     opt.URL,
		"group.id":              opt.GroupId,
		"broker.address.family": "v4",
		"auto.offset.reset":     "latest",
		"enable.auto.commit":    false,
		"max.poll.interval.ms":  60000,
		"debug":                 "consumer",
	})
	if err != nil {
		return nil, cleanup, err
	}
	addCleanup(func() {
		if err := consumer.Close(); err != nil {
			logger.FromContext(ctx).Error("error starting for subscription", opt)
		}
	})

	logger.FromContext(ctx).Info("Created", "Consumer", consumer)
	if err := consumer.SubscribeTopics([]string{opt.TopicName},
		rebalanceCallback); err != nil {
		logger.FromContext(ctx).Error("Error on subscribing to topics", "Error", err)
		return nil, cleanup, err
	}
	return &Subscription{
		consumer: consumer,
		handler:  h,
		opt:      opt,
		msg:      make(chan *kafka.Message),
		cleanup:  cleanup,
	}, cleanup, nil
}

func (s *Subscription) Start(ctx context.Context) error {
	var wgReceiver sync.WaitGroup
	var wgWorker sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wgReceiver.Add(1)
	go s.receiver(ctx, &wgReceiver)

	wgWorker.Add(s.opt.WorkTask)
	for i := 0; i < s.opt.WorkTask; i++ {
		go s.subscribe(ctx, &wgWorker)
	}

	go func() {
		wgReceiver.Wait()
		close(s.msg)
	}()

	wgWorker.Wait()
	return nil
}

func (s *Subscription) receiver(ctx context.Context, wg *sync.WaitGroup) {
	logger.FromContext(ctx).Info("start receiving message")
	defer wg.Done()

	consumed := 0
	paused := false

	for {
		select {
		case <-ctx.Done():
			logger.FromContext(ctx).Info("context canceled, closing receiver")
			return
		default:
			{
				if !paused {
					ev := s.consumer.Poll(s.opt.Poll)
					if ev == nil {
						continue
					}
					switch msg := ev.(type) {
					case *kafka.Message:
						consumed++
						if msg != nil {
							s.msg <- msg
						}

						if consumed >= s.opt.BatchLimit {
							paused = true
						}
					case kafka.PartitionEOF:
						logger.FromContext(ctx).Error("%% Reached", "error", msg)
					case kafka.Error:
						logger.FromContext(ctx).Error("%%", "Error", msg)
					default:
						if msg.String() != "" {
							logger.FromContext(ctx).Error(
								"Ignored", "Error", msg)
						}
					}
				} else {
					if len(s.msg) == 0 {
						consumed = 0
						paused = false
					}
				}
			}
		}
	}
}

func (s *Subscription) doCommit(ctx context.Context) error {
	info, err := s.consumer.Commit()
	if len(info) != 0 {
		logger.FromContext(ctx).Info("Committed", "Topic", *info[len(info)-1].Topic, "Offset",
			info[len(info)-1].Offset.String())
	}
	if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
		logger.FromContext(ctx).Error("Error on committing", "Error", err)
		return nil
	}
	return err
}

func (s *Subscription) subscribe(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.FromContext(ctx).Errorf("recovered from panic: %v", r)
		}
		wg.Done()
	}()

	var msgCount int64

	batchTimeout := 2 * time.Second
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.FromContext(ctx).Info("context canceled, closing")
		case <-timer.C:
			if msgCount != 0 {
				if err := s.doCommit(ctx); err != nil {
					logger.FromContext(ctx).Error("Error on doCommit", err)
				}
				msgCount = 0
			}
			timer.Reset(batchTimeout)
		case msg, ok := <-s.msg:
			if !ok {
				logger.FromContext(ctx).Info("message channel closed, exiting worker")
				return
			}
			if err := s.handler(ctx, msg.Value); err != nil {
				logger.FromContext(ctx).Error("Error on processing message", "Error", err)
				continue
			}
			msgCount++
			if msgCount == int64(s.opt.BatchSize) {
				if err := s.doCommit(ctx); err != nil {
					logger.FromContext(ctx).Error("Error on doCommit", err)
					continue
				}
				msgCount = 0
				timer.Reset(batchTimeout)
			}
		}
	}
}

func rebalanceCallback(c *kafka.Consumer, e kafka.Event) error {
	switch ev := e.(type) {
	case kafka.AssignedPartitions:
		fmt.Printf("%% %s rebalance: %d new partition(s) assigned: %v\n",
			c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)
		err := c.Assign(ev.Partitions)
		if err != nil {
			return err
		}
	case kafka.RevokedPartitions:
		fmt.Printf("%% %s rebalance: %d partition(s) revoked: %v\n",
			c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)

		if c.AssignmentLost() {
			slog.Info("Assignment lost involuntarily, commit may fail")
		}

		// Since enable.auto.commit is unset, we need to commit offsets manually
		// before the partition is revoked.
		commitedOffsets, err := c.Commit()

		if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
			slog.Error("Failed to commit offsets", slog.Any("error", err))
			return err
		}
		slog.Info("Committed offsets to Kafka", slog.Any("commitedOffsets", commitedOffsets))
	default:
		slog.Error("Unexpected event type", slog.Any("event", e))
	}

	return nil

}
