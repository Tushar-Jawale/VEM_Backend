package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/akashtripathi12/TBO_Backend/internal/utils"
	"github.com/hibiken/asynq"
)

// HandleEmailTask processes the email delivery task
func HandleEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("📨 [WORKER] Processing email task for: %s", p.To)

	// Use our existing email utility
	// Note: SendEmail expects a slice of strings for 'to'
	err := utils.SendEmail([]string{p.To}, p.Subject, p.Body)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
