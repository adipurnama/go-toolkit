package pubsubkit

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// NewPubSubSubscription returns new PubSub topic subscriber in 5s timeout.
func NewPubSubSubscription(
	client *pubsub.Client,
	subID string,
	cfg pubsub.SubscriptionConfig,
	opts ...Option,
) (*pubsub.Subscription, error) {
	opt := newDefaultOptions()

	for _, o := range opts {
		o(opt)
	}

	sub := client.Subscription(subID)

	connectCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	if opt.checkExists {
		ok, err := sub.Exists(connectCtx)
		if err != nil {
			return nil, errors.Wrapf(err, "pubsubkit: failed to check subscription %s existence", subID)
		}

		if !ok {
			if !opt.autoCreate {
				return nil, errors.Wrapf(ErrSubscriptionNotFound, "pubsubkit: failed create subscription %s", subID)
			}

			sub, err = client.CreateSubscription(connectCtx, subID, cfg)
			if err != nil {
				return nil, errors.Wrapf(err, "pubsubkit: failed create pubsub subscription %s", subID)
			}
		}
	}

	return sub, nil
}
