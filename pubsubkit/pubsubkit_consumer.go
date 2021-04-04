package pubsubkit

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// NewPubSubConsumer returns new PubSub topic subscriber in 5s timeout.
func NewPubSubConsumer(
	client *pubsub.Client,
	topic *pubsub.Topic,
	subID string,
	autoCreateSubscription bool,
) (*pubsub.Subscription, error) {
	sub := client.Subscription(subID)

	connectCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	ok, err := sub.Exists(connectCtx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check subscription existence")
	}

	if !ok {
		if !autoCreateSubscription {
			return nil, errors.Wrapf(ErrSubscriptionNotFound, "failed create subscription %s", subID)
		}

		sub, err = client.CreateSubscription(connectCtx, subID, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed create pubsub subscription")
		}
	}

	return sub, nil
}
