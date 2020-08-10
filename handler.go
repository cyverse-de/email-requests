package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cyverse-de/logcabin"

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logcabin.Error.Printf("unable to read error response body: %s", err.Error())
		return
	}

	// Unmarshal the response body.
	var errorResponse cyverseEmailErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		logcabin.Error.Printf("unable to parse error response body: %s", err.Error())
	}

	// Log the response body.
	logcabin.Error.Printf("cyverse-email returned an error: %s", errorResponse.Message)
}

// HandleMessage handles a single incoming AMQP delivery. A communication error with cyverse-email probably
// means that cyverse-email is down. The service aborts in that case because no useful work can be done if
// cyverse-email is down.
func (h *Handler) HandleMessage(delivery amqp.Delivery) error {
	// Forward the request to the cyverse-email service.
	resp, err := http.Post(h.cyverseEmailBaseURL, "application/json", bytes.NewReader(delivery.Body))
	if err != nil {
		logcabin.Error.Fatalf("unable to communicate with cyverse-email: %s", err.Error())
	}
	defer resp.Body.Close()

	// Check the status of the response. We don't requeue the message because it'll probably fail again.
	if resp.StatusCode >= http.StatusBadRequest {
		h.logErrorResponse(resp)
	}

	// Acknowledge the message.
	err = delivery.Ack(false)
	if err != nil {
		logcabin.Error.Printf("unable to acknowledge message: %s", err.Error())
	}

	return nil
}
