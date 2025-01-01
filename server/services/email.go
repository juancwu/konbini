package services

import (
	"bytes"
	"context"
	"fmt"
	"konbini/server/config"
	"konbini/server/views"

	"github.com/resend/resend-go/v2"
)

// SendEmail sends an email via the Resend Client. This is the base function and
// ideally not used directly but instead as the only step where an email is sent.
func SendEmail(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	c, err := config.Global()
	if err != nil {
		return nil, err
	}
	client := resend.NewClient(c.GetResendApiKey())
	sent, err := client.Emails.SendWithContext(ctx, params)
	return sent, err
}

// SendVerificationEmail sends an email verification for users to verify their email.
func SendVerificationEmail(ctx context.Context, to string, token string) (*resend.SendEmailResponse, error) {
	c, err := config.Global()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/auth/email/verify?token=%s", c.GetBackendUrl(), token)
	component := views.VerificationEmail(url)
	var buffer bytes.Buffer
	err = component.Render(ctx, &buffer)
	if err != nil {
		return nil, err
	}

	params := &resend.SendEmailRequest{
		From:    c.GetVerifyEmailAddress(),
		To:      []string{to},
		Subject: "Verify Your Email",
		Html:    buffer.String(),
		Text: fmt.Sprintf(
			`Thanks for using Konbini!

Please verify your email by opening the following link in a browser:

%s`,
			url,
		),
	}

	res, err := SendEmail(ctx, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}
