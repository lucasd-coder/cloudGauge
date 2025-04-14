package scheduledtaskrunner

import "golang.org/x/net/context"

func Start(ctx context.Context) error {
	if err := StartDispatchScheduler(ctx); err != nil {
		return err
	}

	return nil
}
