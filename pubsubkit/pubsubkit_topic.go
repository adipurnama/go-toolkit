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
	opts ...Option,
) (*pubsub.Topic, error) {
	opt := newDefaultOptions()

	for _, o := range opts {
		o(opt)
	}

	topic := client.Topic(topicID)

	connectCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	if opt.checkExists {
		ok, err := topic.Exists(connectCtx)
		if err != nil {
			return nil, errors.Wrapf(err, "pubsubkit: failed to check topic %s existence", topicID)
		}

		if !ok {
			if !opt.autoCreate {
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
	}

	return topic, nil
}
