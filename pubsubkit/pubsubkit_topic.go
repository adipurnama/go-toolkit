package pubsubkit

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// NewPubSubTopic returns new PubSub topic publisher in 5s timeout.
func NewPubSubTopic(
	client *pubsub.Client,
	topicID string,
	cfg *pubsub.TopicConfig,
) (*pubsub.Topic, error) {
	return newPubSubTopicWithAutocreateOption(client, topicID, cfg, false)
}

// NewPubSubTopicAutocreate returns new PubSub topic publisher in 5s timeout.
// it will create new topic if it doesn't exist yet on the server.
func NewPubSubTopicAutocreate(
	client *pubsub.Client,
	topicID string,
	cfg *pubsub.TopicConfig,
) (*pubsub.Topic, error) {
	return newPubSubTopicWithAutocreateOption(client, topicID, cfg, true)
}

func newPubSubTopicWithAutocreateOption(
	client *pubsub.Client,
	topicID string,
	cfg *pubsub.TopicConfig,
	autoCreateTopic bool,
) (*pubsub.Topic, error) {
	topic := client.Topic(topicID)

	connectCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	ok, err := topic.Exists(connectCtx)
	if err != nil {
		return nil, errors.Wrapf(err, "pubsubkit: failed to check topic %s existence", topicID)
	}

	if !ok {
		if !autoCreateTopic {
			return nil, errors.Wrapf(ErrTopicNotFound, "pubsubkit: failed creating publisher for topic %s", topicID)
		}

		if cfg != nil {
			topic, err = client.CreateTopicWithConfig(connectCtx, topicID, cfg)
		} else {
			topic, err = client.CreateTopic(connectCtx, topicID)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "pubsubkit: failed to create pubsub topic %s", topicID)
		}
	}

	return topic, nil
}
