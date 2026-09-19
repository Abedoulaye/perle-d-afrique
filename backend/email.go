package main

import (
	"fmt"
	"os"
	"github.com/resend/resend-go/v3"
)

func sendVerificationEmail(toEmail, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	verifyURL := fmt.Sprintf("%s/verify?token=%s", os.Getenv("CORS_ORIGIN"), token)

	params := &resend.SendEmailRequest {
		From: "onboarding@resend.dev", // replace with domain later
		To: []string{toEmail},
		Subject: "Verify your email - Diariata's Fabrics",
		Html: fmt.Sprintf(`
		<h1>Welcome to Diariata's Fabrics!</h1>
		<p>Click the link below to verify your email:</p>
		<a href="%s">Verify Email </a>
		<p>This link expires in 24 hours.</p>
		`, verifyURL)
	}

	_, err := client.Emails.Send(params)
	return err
}