package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/deepakgudla/BookVault/internal/config"
	"github.com/deepakgudla/BookVault/internal/models"
	"github.com/deepakgudla/BookVault/internal/providers"
	"github.com/deepakgudla/BookVault/notifications"
)

func main() {
	log.Println("notification service has started...")

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// email notifier
	emailConfig := &notifications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}

	emailNotifier := notifications.NewEmailNotifier(emailConfig)

	// AWS config for SQS
	awsConfig, err := providers.CreateAWSConfig(ctx, cfg.AWS.S3Endpoint, cfg.AWS.Region)
	if err != nil {
		log.Fatalf("failed to create AWS config: %v", err)
	}

	// SQS subscriber
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig: awsConfig,
	}, logger)

	if err != nil {
		log.Fatalf("failed to create subscrier: %v", err)
	}

	messages, err := subscriber.Subscribe(ctx, cfg.AWS.EventQueueName)
	if err != nil {
		if err := subscriber.Close(); err != nil {
			log.Printf("failed to close subscriber: %v", err)
		}
		log.Fatalf("failed to subscribe to queue: %v", err)
	}

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("notification server has started, waiting for messages")

	for {
		select {
		case msg := <-messages:
			if err := processMessage(msg, emailNotifier); err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case <-sigCh:
			log.Println("shutting down notification service...")

			if err := subscriber.Close(); err != nil {
				log.Printf("failed to close subscriber: %v", err)
			}

			return
		}
	}
}

func processMessage(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	eventType := msg.Metadata.Get("event_type")
	switch eventType {
	case "USER_LOGGED_IN":
		return handleUserLoggedIn(msg, emailNotifier)
	default:
		log.Printf("unknown event type: %s", eventType)
		return nil
	}
}

func handleUserLoggedIn(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	var user models.User
	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return err
	}

	userName := user.FirstName + " " + user.LastName
	if userName == " " {
		userName = "User"
	}

	log.Printf("sending login notification to %s", user.Email)

	return emailNotifier.SendLoginNotification(user.Email, userName)
}
