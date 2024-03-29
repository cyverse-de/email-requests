package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/streadway/amqp"
)

// Handler is the AMQP message handler for this service.
type Handler struct {
	cyverseEmailBaseURL string
}

// NewHandler returns a new AMQP message handler.
func NewHandler(cyverseEmailBaseURL string) *Handler {
	return &Handler{cyverseEmailBaseURL: cyverseEmailBaseURL}
}

// cyverseEmailErrorResponse is an error response body returned by the cyverse-email service.
type cyverseEmailErrorResponse struct {
	Message string `json:"message"`
}

// logErrorResponse logs the error message in a response from cyverse-email.
func (*Handler) logErrorResponse(resp *http.Response) {

	// Slurp the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("unable to read error response body: %s", err.Error())
		return
	}

	// Unmarshal the response body.
	var errorResponse cyverseEmailErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		log.Errorf("unable to parse error response body: %s", err.Error())
	}

	// Log the response body.
	log.Errorf("cyverse-email returned an error: %s", errorResponse.Message)
}

// HandleMessage handles a single incoming AMQP delivery. A communication error with cyverse-email probably
// means that cyverse-email is down. The service aborts in that case because no useful work can be done if
// cyverse-email is down.
func (h *Handler) HandleMessage(ctx context.Context, delivery amqp.Delivery) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.cyverseEmailBaseURL, bytes.NewReader(delivery.Body))
	if err != nil {
		log.Fatalf("unable to communicate with cyverse-email: %s", err.Error())
	}

	req.Header.Set("content-type", "application/json")

	// Forward the request to the cyverse-email service.
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("unable to communicate with cyverse-email: %s", err.Error())
	}
	defer resp.Body.Close()

	// Check the status of the response. We don't requeue the message because it'll probably fail again.
	if resp.StatusCode >= http.StatusBadRequest {
		h.logErrorResponse(resp)
	}

	// Acknowledge the message.
	err = delivery.Ack(false)
	if err != nil {
		log.Errorf("unable to acknowledge message: %s", err.Error())
	}

	return nil
}
