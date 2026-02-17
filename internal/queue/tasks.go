package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

// Task Types
const (
	TypeEmailDelivery = "email:deliver"
)

// EmailPayload is the data passed to the task
type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// NewEmailTask creates a task for email delivery
func NewEmailTask(to string, subject string, body string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}
