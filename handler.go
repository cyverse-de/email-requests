package main

import "github.com/streadway/amqp"

// Handler is the AMQP message handler for this service.
type Handler struct {
	cyverseEmailBaseURL string
}

// NewHandler returns a new AMQP message handler.
func NewHandler(cyverseEmailBaseURL string) *Handler {
	return &Handler{cyverseEmailBaseURL: cyverseEmailBaseURL}
}

// HandleMessage handles a single incoming AMQP delivery.
func (h *Handler) HandleMessage(delivery amqp.Delivery) error {
	return nil
}
