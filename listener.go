package main

import (
	"github.com/cyverse-de/logcabin"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gopkg.in/cyverse-de/messaging.v7"
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
func (l *Listener) handleMessage(delivery amqp.Delivery) {
	err := l.handler.HandleMessage(delivery)
	if err != nil {
		logcabin.Error.Printf("Error occurred while handling message: %s", err.Error())
	}
}

// Listen listens for and handles incoming AMQP messages.
func (l *Listener) Listen() {
	// Listen for incoming messages.
	l.amqpClient.AddConsumer(
		l.amqpSettings.ExchangeName,
		l.amqpSettings.ExchangeType,
		queueName,
		messaging.EmailRequestPublishingKey,
		l.handleMessage,
		100,
	)
}
