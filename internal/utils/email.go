package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends an email using Gmail SMTP
// It reads SMTP_EMAIL and SMTP_PASS from environment variables
// If credentials are missing, it logs the email to the console (useful for dev)
func SendEmail(to []string, subject string, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASS")

	// Development mode: Log if credentials missing
	if from == "" || password == "" {
		log.Printf("⚠️ [DEV MODE] Email to %v\nSubject: %s\nBody (truncated): %s...\n", to, subject, body[:min(len(body), 50)])
		log.Println("ℹ️ Set SMTP_EMAIL and SMTP_PASS to send real emails.")
		return nil
	}

	// SMTP Server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message construction
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", to[0], subject, body))

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("📧 Email sent to %v", to)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
