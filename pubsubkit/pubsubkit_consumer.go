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
		return nil, errors.Wrapf(err, "pubsubkit: failed to check subscription %s existence", subID)
	}

	if !ok {
		if !autoCreateSubscription {
			return nil, errors.Wrapf(ErrSubscriptionNotFound, "pubsubkit: failed create subscription %s", subID)
		}

		sub, err = client.CreateSubscription(connectCtx, subID, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "pubsubkit: failed create pubsub subscription %s", subID)
		}
	}

	return sub, nil
}
