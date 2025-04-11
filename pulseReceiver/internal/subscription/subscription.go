package subscription

import "context"

func Start(ctx context.Context) error {
	if err := PulseReceiverEvent(ctx); err != nil {
		return err
	}

	return nil
}
