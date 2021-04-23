// Package pubsubkit provides helper to interact with GCP PubSub
package pubsubkit

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

var (
	// ErrSubscriptionNotFound ...
	ErrSubscriptionNotFound = errors.New("pubsub subscription doesn't exists")
	// ErrTopicNotFound ...
	ErrTopicNotFound = errors.New("pubsub topic doesn't exists")
)

const (
	connectTimeout = 5 * time.Second
)

// NewPubSubClient returns new PubSub client in 5s timeout.
func NewPubSubClient(projectID string, opts ...option.ClientOption) (*pubsub.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "pubsubkit: failed to create pubsub client")
	}

	return client, nil
}
