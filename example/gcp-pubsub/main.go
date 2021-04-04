package main

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/pubsubkit"
	"github.com/adipurnama/go-toolkit/runtimekit"
)

const (
	projectID = "test-project"
	topicID   = "topic.example"
	// topicDltID     = "topic.example.dlt".
	subscriptionID = "test-sub_" + topicID

	isProdMode = false
	appName    = "sample-pubsub-app"

	publishInterval = 2 * time.Second
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
	topicPublisher, err := pubsubkit.NewPubSubPublisher(client, topicID, true)
	if err != nil {
		log.FromCtx(appCtx).Error(err, "failed to create topic publisher")
		return
	}

	log.FromCtx(appCtx).Info("topic publisher created")

	// GCP PubSub Subscriber setup
	sub, err := pubsubkit.NewPubSubConsumer(client, topicPublisher, subscriptionID, true)
	if err != nil {
		log.FromCtx(appCtx).Error(err, "failed to create topic subscriber")
		return
	}

	log.FromCtx(appCtx).Info("subscriptions created")

	go readMessages(appCtx, sub)

	// publish message
	log.Println("start publishing messages...")

	for {
		select {
		case <-time.After(publishInterval):
			publishMessages(topicPublisher)
		case <-appCtx.Done():
			log.FromCtx(appCtx).Info("exit signal received. stopping publisher...")
			topicPublisher.Stop()
			log.Println("Bye")

			return
		}
	}
}

func readMessages(appCtx context.Context, sub *pubsub.Subscription) {
	err := sub.Receive(appCtx, func(ctx context.Context, msg *pubsub.Message) {
		log.FromCtx(ctx).Info("got message", "data", msg.Data, "msg_id", msg.ID)
		msg.Ack()
	})
	if err != nil {
		log.FromCtx(appCtx).Error(err, "error while receiving message")
	}
}

func publishMessages(publisher *pubsub.Topic) {
	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	msg := &pubsub.Message{
		Data: []byte("Hello, world"),
	}
	result := publisher.Publish(ctx, msg)

	msgID, err := result.Get(ctx)
	if err != nil {
		log.FromCtx(ctx).Error(err, "publishing message failed")
		return
	}

	log.FromCtx(ctx).Info("message published", "msg_id", msgID)
}
