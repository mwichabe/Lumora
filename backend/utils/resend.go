package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"lumora/backend/config"
)

// Email delivery through the Resend HTTP API.
//
// This exists because the production API runs on a host that blocks outbound
// SMTP ports: the Gmail SMTP path works on a laptop and silently times out once
// deployed. Resend is plain HTTPS on 443, which is never blocked.

// resendEndpoint is a var so tests can point it at a local server.
var resendEndpoint = "https://api.resend.com/emails"

// resendTimeout bounds a single send. The call runs in a goroutine, so this
// only stops a hung connection from leaking — it never delays a response.
const resendTimeout = 15 * time.Second

var resendClient = &http.Client{Timeout: resendTimeout}

type resendPayload struct {
	From    string            `json:"from"`
	To      []string          `json:"to"`
	Subject string            `json:"subject"`
	HTML    string            `json:"html"`
	Text    string            `json:"text"`
	Headers map[string]string `json:"headers,omitempty"`
}

// sendViaResend posts one message to Resend. A non-2xx reply is returned as an
// error carrying Resend's own explanation (unverified domain, bad key, ...),
// which is what you want to see in the server logs when mail goes missing.
func sendViaResend(cfg config.Config, toEmail, subject, plain, html string) error {
	body, err := json.Marshal(resendPayload{
		From:    cfg.ResendFrom,
		To:      []string{toEmail},
		Subject: subject,
		HTML:    html,
		Text:    plain,
		// Signal that this is an automated, no-reply message.
		Headers: map[string]string{
			"Auto-Submitted":           "auto-generated",
			"X-Auto-Response-Suppress": "All",
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := resendClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		detail, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return fmt.Errorf("resend responded %d: %s", res.StatusCode, bytes.TrimSpace(detail))
	}
	return nil
}
