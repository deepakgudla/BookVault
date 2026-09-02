package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deepakgudla/bookvault/internal/providers"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"

	_ "github.com/aws/smithy-go/endpoints"
	appconfig "github.com/deepakgudla/bookvault/internal/config"
)

// EventPublisher publishes application events to an SQS queue.
type EventPublisher struct {
	publisher message.Publisher
	queueName string
}

// Publish serializes and publishes an event with its metadata.
func (ep *EventPublisher) Publish(eventType string, payload interface{}, metadata map[string]string) error {

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)

	msg.Metadata.Set("event_type", eventType)

	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}

	return ep.publisher.Publish(ep.queueName, msg)
}

// Close releases the underlying event publisher.
func (ep *EventPublisher) Close() error {
	return ep.publisher.Close()
}

// NewEventPublisher creates an SQS-backed event publisher.
func NewEventPublisher(ctx context.Context, cfg *appconfig.AWSConfig) (*EventPublisher, error) {
	logger := watermill.NewStdLogger(false, false)

	awsConfig, err := providers.CreateAWSConfig(ctx, cfg.S3Endpoint, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to crete AWS config: %w", err)
	}

	publisherConfig := sqs.PublisherConfig{
		AWSConfig: awsConfig,
		Marshaler: nil,
	}

	publisher, err := sqs.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	fmt.Printf("EventQueueName = '%s'\n", cfg.EventQueueName)

	return &EventPublisher{
		publisher: publisher,
		queueName: cfg.EventQueueName,
	}, nil
}
