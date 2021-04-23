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
) (*pubsub.Subscription, error) {
	return newPubsubConsumerWithAutoCreateOption(client, subID, cfg, false)
}

// NewPubSubSubscriptionAutocreate returns new PubSub topic subscription in 5s timeout
// when subscription is not found, it'll try to create it instead return error directly.
func NewPubSubSubscriptionAutocreate(
	client *pubsub.Client,
	subID string,
	cfg pubsub.SubscriptionConfig,
) (*pubsub.Subscription, error) {
	return newPubsubConsumerWithAutoCreateOption(client, subID, cfg, true)
}

func newPubsubConsumerWithAutoCreateOption(
	client *pubsub.Client,
	subID string,
	cfg pubsub.SubscriptionConfig,
	autoCreate bool,
) (*pubsub.Subscription, error) {
	sub := client.Subscription(subID)

	connectCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	ok, err := sub.Exists(connectCtx)
	if err != nil {
		return nil, errors.Wrapf(err, "pubsubkit: failed to check subscription %s existence", subID)
	}

	if !ok {
		if !autoCreate {
			return nil, errors.Wrapf(ErrSubscriptionNotFound, "pubsubkit: failed create subscription %s", subID)
		}

		sub, err = client.CreateSubscription(connectCtx, subID, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "pubsubkit: failed create pubsub subscription %s", subID)
		}
	}

	return sub, nil
}
