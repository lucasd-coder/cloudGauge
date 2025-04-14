package jobscheduled

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
)

type JobScheduler struct {
	handler     func(ctx context.Context)
	scheduled   gocron.Scheduler
	cleanup     func()
	durationJob time.Duration
}

func NewJobScheduler(ctx context.Context,
	h func(ctx context.Context),
	durationJob time.Duration,
) (_ *JobScheduler, cleanup func(), err error) {
	addCleanup := func(f func()) {
		old := cleanup
		cleanup = func() { old(); f() }
	}

	log := logger.FromContext(ctx)

	defer func() {
		if err != nil {
			cleanup()
			cleanup = nil
		}
	}()

	cleanup = func() {}

	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, cleanup, err
	}
	addCleanup(func() {
		if shutdownErr := s.Shutdown(); shutdownErr != nil {
			log.Error("error shutting down scheduler", "error", shutdownErr)
		}
	})

	log.Info("Scheduler created successfully")
	return &JobScheduler{
		handler:     h,
		cleanup:     cleanup,
		scheduled:   s,
		durationJob: durationJob,
	}, cleanup, nil
}

func (s *JobScheduler) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log := logger.FromContext(ctx)

	job, err := s.scheduled.NewJob(
		gocron.DurationJob(s.durationJob), gocron.NewTask(
			s.handler),
		gocron.WithContext(ctx))
	if err != nil {
		return err
	}
	log.Info("Task started", "job", job.Name())

	s.scheduled.Start()

	<-ctx.Done()
	log.Info("Scheduler context canceled")

	return nil
}
