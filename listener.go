package main

import (
	"context"

	"github.com/cyverse-de/messaging/v9"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const queueName = "email_requests"

// Listener is the AMQP message listener for this service.
type Listener struct {
	amqpClient   *messaging.Client
	amqpSettings *AMQPSettings
	handler      *Handler
}

// NewListener returns a new AMQP message listener.
func NewListener(handler *Handler, amqpSettings *AMQPSettings) (*Listener, error) {
	wrapMsg := "unable to create the message listener"

	// Create the AMQP client.
	amqpClient, err := messaging.NewClient(amqpSettings.URI, false)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Build and return the listener.
	listener := &Listener{
		amqpClient:   amqpClient,
		amqpSettings: amqpSettings,
		handler:      handler,
	}
	return listener, nil
}

// handleMessage handles a single incoming AMQP message.
func (l *Listener) handleMessage(ctx context.Context, delivery amqp.Delivery) {
	err := l.handler.HandleMessage(ctx, delivery)
	if err != nil {
		log.Errorf("error occurred while handling message: %s", err.Error())
	}
}

// Listen listens for and handles incoming AMQP messages.
func (l *Listener) Listen() {
	go l.amqpClient.Listen()

	// Add a consumer.
	l.amqpClient.AddConsumer(
		l.amqpSettings.ExchangeName,
		l.amqpSettings.ExchangeType,
		queueName,
		messaging.EmailRequestPublishingKey,
		l.handleMessage,
		100,
	)
}
