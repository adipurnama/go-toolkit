package pubsubkit

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"

	"github.com/adipurnama/go-toolkit/log"
)

type (
	// WorkerHandlerFunc handles single message received from *pubsub.Subscription.
	// Your handler should be idempotent since gcp PubSub might send message more than once
	// or message arrived out of order.
	WorkerHandlerFunc func(ctx context.Context, msg Message) error

	// Message wraps *pubsub.Message without the Nack() & Ack() handler
	// it is designed as parameter to `WorkerHandlerFunc`.
	Message interface {
		// ID identifies this message.
		// This ID is assigned by the server and is populated for Messages obtained from a subscription.
		// This field is read-only.
		ID() string

		// Data is the actual data in the message.
		Data() []byte

		// Attributes represents the key-value pairs the current message
		// is labelled with.
		Attributes() map[string]string

		// The time at which the message was published.
		// This is populated by the server for Messages obtained from a subscription.
		// This field is read-only.
		PublishTime() time.Time

		// DeliveryAttempt is the number of times a message has been delivered.
		// This is part of the dead lettering feature that forwards messages that
		// fail to be processed (from nack/ack deadline timeout) to a dead letter topic.
		// If dead lettering is enabled, this will be set on all attempts, starting
		// with value 1. Otherwise, the value will be nil.
		// This field is read-only.
		DeliveryAttempt() int

		// DLTSupported defines wether message's topic has DLT nor not
		DLTSupported() bool

		// OrderingKey identifies related messages for which publish order should
		// be respected. If empty string is used, message will be sent unordered.
		OrderingKey() string
	}

	msgWrapper struct {
		msg *pubsub.Message
	}
)

// ID is Message interface impl.
func (w *msgWrapper) ID() string {
	return w.msg.ID
}

// Data is Message interface impl.
func (w *msgWrapper) Data() []byte {
	return w.msg.Data
}

// Attributes is Message interface impl.
func (w *msgWrapper) Attributes() map[string]string {
	return w.msg.Attributes
}

// PublishTime is Message interface impl.
func (w *msgWrapper) PublishTime() time.Time {
	return w.msg.PublishTime
}

// DeliveryAttempt is Message interface impl.
func (w *msgWrapper) DeliveryAttempt() (count int) {
	if w.msg.DeliveryAttempt != nil {
		return *w.msg.DeliveryAttempt
	}

	return count
}

// OrderingKey is Message interface impl.
func (w *msgWrapper) OrderingKey() string {
	return w.msg.OrderingKey
}

// DLTSupported is Message interface impl.
func (w *msgWrapper) DLTSupported() bool {
	return w.msg.DeliveryAttempt != nil
}

var (
	// ErrInvalidSubscription returned when try to receive message from subscription `nil`.
	ErrInvalidSubscription = errors.New("pubsubkit: subscription cannot be nil")

	// ErrInvalidSubscriptionHandler returned when try to receive message from subscription using nil handler.
	ErrInvalidSubscriptionHandler = errors.New("pubsubkit: handler cannot be nil")
)

/*
ReceiveSubscription blocks to receive messages from pubsub subscription
Call with goroutine if you'd like to do something else in the meantime.

	go func() {
		if err := pubsubkit.ReceiveSubscription(...); err != nil {
		  // handle error
		}
	}()

It will `Nack()` message when handler returns error & DLT found
`Ack()` when handler is success, or error with DLT not found
it also logs the process using `toolkit/log` package.
*/
func ReceiveSubscription(
	ctx context.Context,
	sub *pubsub.Subscription,
	handler WorkerHandlerFunc,
	opts ...Option,
) (err error) {
	opt := newDefaultOptions()

	for _, o := range opts {
		o(opt)
	}

	if sub == nil {
		return errors.WithStack(ErrInvalidSubscription)
	}

	if handler == nil {
		return errors.WithStack(ErrInvalidSubscriptionHandler)
	}

	var cfg pubsub.SubscriptionConfig

	if opt.checkExists {
		cfg, err = sub.Config(ctx)
		if err != nil {
			return errors.Wrap(err, "pubsubkit: get subscription config failed")
		}
	}

	log.FromCtx(ctx).Info("pubsub worker started. listening messages...", "subscription", sub.String(), "config", cfg)

	recErr := sub.Receive(ctx, func(wCtx context.Context, msg *pubsub.Message) {
		logFields := []interface{}{
			"msg", string(msg.Data),
			"msg_id", msg.ID,
			"worker_id", sub.ID(),
			"subscription", sub.String(),
		}

		attempt := 0
		if msg.DeliveryAttempt != nil {
			attempt = *msg.DeliveryAttempt
		}

		logFields = append(logFields, "delivery_attempt", attempt)

		err := handler(wCtx, &msgWrapper{msg})
		if err == nil {
			msg.Ack()
			log.FromCtx(ctx).Info("Message successfully processed & ACK'ed.", logFields...)
			return
		}

		msg.Nack()

		if opt.checkExists && cfg.DeadLetterPolicy == nil {
			log.FromCtx(ctx).Error(err, "Processing message failed. No DLTPolicy found. Message NOT ACK'ed.", logFields...)
			return
		}

		log.FromCtx(ctx).Error(err, "Processing message failed. Message NOT ACK'ed.", logFields...)
	})

	return errors.Wrap(recErr, "pubsubkit: error while receiving subscription messages")
}
