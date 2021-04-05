package pubsubkit

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// NewPubSubPublisher returns new PubSub topic publisher in 5s timeout.
func NewPubSubPublisher(
	client *pubsub.Client,
	topicID string,
	autoCreateTopic bool,
) (*pubsub.Topic, error) {
	topicPublisher := client.Topic(topicID)

	connectCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	ok, err := topicPublisher.Exists(connectCtx)
	if err != nil {
		return nil, errors.Wrapf(err, "pubsubkit: failed to check topic %s existence", topicID)
	}

	if !ok {
		if !autoCreateTopic {
			return nil, errors.Wrapf(ErrTopicNotFound, "pubsubkit: failed creating publisher for topic %s", topicID)
		}

		topicPublisher, err = client.CreateTopic(connectCtx, topicID)
		if err != nil {
			return nil, errors.Wrapf(err, "pubsubkit: failed to create pubsub topic %s", topicID)
		}
	}

	return topicPublisher, nil
}
