package main

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	shortuuid "github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/pubsubkit"
	"github.com/adipurnama/go-toolkit/runtimekit"
)

const (
	projectID           = "test-project"
	topicID             = "topic.example"
	errorTopicID        = "topic.example.error"
	subscriptionID      = "topic.example_subscription"
	errorSubscriptionID = "topic.example.error_subscription"

	isProdMode = false
	appName    = "sample-pubsub-app"

	publishInterval = 5 * time.Second
	publishTimeout  = 3 * time.Second
)

func main() {
	appCtx, cancel := runtimekit.NewRuntimeContext()
	defer cancel()

	// setup logging
	if isProdMode {
		// production mode - json
		_ = log.NewLogger(log.LevelDebug, appName, nil, nil, "additionalKey1", "additional_value1").Set()
	} else {
		// development mode - logfmt
		_ = log.NewDevLogger(nil, nil).Set()
	}

	// create pubsub client
	client, err := pubsubkit.NewPubSubClient(projectID)
	if err != nil {
		log.FromCtx(appCtx).Error(err, "failed create pubsub client")
		return
	}

	log.FromCtx(appCtx).Info("connected to pubsub server")

	// Setup GCP PubSub topic publisher
	topic, err := pubsubkit.NewPubSubTopic(client, topicID, &pubsub.TopicConfig{}, pubsubkit.WithAutoCreate())
	if err != nil {
		log.FromCtx(appCtx).Error(err, "failed to create topic")
		return
	}

	errorTopic, err := pubsubkit.NewPubSubTopic(client, errorTopicID, &pubsub.TopicConfig{}, pubsubkit.WithAutoCreate())
	if err != nil {
		log.FromCtx(appCtx).Error(err, "failed to create error topic")
		return
	}

	log.FromCtx(appCtx).Info("pubsub topics created")

	// GCP PubSub Subscriber setup
	sub, err := pubsubkit.NewPubSubSubscription(client, subscriptionID, pubsub.SubscriptionConfig{
		Topic: topic,
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			// DLT should be fullyQualifiedProjectName ID, see: https://cloud.google.com/pubsub/docs/dead-letter-topics
			// you can't simply pass string of `errorTopicID`
			//
			// MaxDeliveryAttempts range from 5-100
			DeadLetterTopic: errorTopic.String(), MaxDeliveryAttempts: 5,
		},
		RetryPolicy: &pubsub.RetryPolicy{
			MinimumBackoff: 1 * time.Second,
			MaximumBackoff: 3 * time.Second,
		},
	}, pubsubkit.WithAutoCreate())
	if err != nil {
		log.FromCtx(appCtx).Error(err, "failed to create topic subscriber")
		return
	}

	errorSub, err := pubsubkit.NewPubSubSubscription(client, errorSubscriptionID, pubsub.SubscriptionConfig{
		Topic: errorTopic,
		RetryPolicy: &pubsub.RetryPolicy{
			MinimumBackoff: 1 * time.Second,
			MaximumBackoff: 3 * time.Second,
		},
	}, pubsubkit.WithAutoCreate())
	if err != nil {
		log.FromCtx(appCtx).Error(err, "failed to create topic subscriber")
		return
	}

	log.FromCtx(appCtx).Info("subscriptions created")

	go func() {
		if err := pubsubkit.ReceiveSubscription(appCtx, sub, readMessages); err != nil {
			log.FromCtx(appCtx).Error(err, "error reading subscription")
		}
	}()

	go func() {
		if err := pubsubkit.ReceiveSubscription(appCtx, errorSub, readErrorMessages); err != nil {
			log.FromCtx(appCtx).Error(err, "error reading error subscription")
		}
	}()

	// publish message
	log.Println("start publishing messages...")

	for i := 0; i < 5; i++ {
		publishMessages(topic)
	}

	<-appCtx.Done()
	log.FromCtx(appCtx).Info("exit signal received. stopping publisher...")
	topic.Stop()
	log.Println("Bye")
}

func readMessages(ctx context.Context, msg pubsubkit.Message) error {
	msgID, err := strconv.Atoi(msg.ID())
	if err != nil {
		return errors.Wrapf(err, "subscriber: failed to parse msg_id=%s", msg.ID())
	}

	// simulates processing failures
	// for each msgID, if it ID%5 == 0, returns error
	// based on subscription settings, it will be retry'ed up to 5 times,
	// then enter the DLT topic (readErrorMessages)
	err = errors.New("failed processing message for failedID")

	failedID := 5
	if msgID%failedID == 0 {
		return errors.Wrapf(err, "subscriber: failed processing msg_id=%s", msg.ID())
	}

	// simulates success processing
	log.FromCtx(ctx).Info("subscriber: got message", "data", msg.Data(), "msg_id", msg.ID())

	return nil
}

func readErrorMessages(ctx context.Context, msg pubsubkit.Message) error {
	attempt := msg.DeliveryAttempt()

	err := errors.New("error processing error subscription")

	msgID, _ := strconv.Atoi(msg.ID())
	if msgID%3 == 0 {
		log.FromCtx(ctx).Error(err, "subscriber: got message",
			"data", msg.Data(),
			"msg_id", msg.ID(),
			"delivery_attempt", attempt)

		return err
	}

	log.FromCtx(ctx).Warn("subscriber: got message",
		"data", msg.Data(),
		"msg_id", msg.ID(),
		"delivery_attempt", attempt)

	return nil
}

func publishMessages(publisher *pubsub.Topic) {
	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	random := shortuuid.New()

	msg := &pubsub.Message{
		Data: []byte("Hello, world " + random),
	}
	result := publisher.Publish(ctx, msg)

	msgID, err := result.Get(ctx)
	if err != nil {
		log.FromCtx(ctx).Error(err, "publishing message failed")
		return
	}

	log.FromCtx(ctx).Info("publisher: message published", "msg_id", msgID)
}
