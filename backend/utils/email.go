package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"lumora/backend/config"
)

// SendEmail is the core sender: builds a MIME multipart message and sends it via
// Gmail SMTP (STARTTLS on port 587) using the App Password from config. Safe to
// call in a goroutine — failures are logged, never fatal, and if SMTP isn't
// configured it simply no-ops.
func SendEmail(cfg config.Config, toEmail, subject, plain, html string) error {
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		log.Printf("[email] SMTP not configured — skipping '%s' to %s", subject, toEmail)
		return nil
	}
	msg := buildMIME(cfg, toEmail, subject, plain, html)
	addr := cfg.SMTPHost + ":" + cfg.SMTPPort
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	if err := smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{toEmail}, []byte(msg)); err != nil {
		log.Printf("[email] failed to send '%s' to %s: %v", subject, toEmail, err)
		return err
	}
	log.Printf("[email] sent '%s' to %s", subject, toEmail)
	return nil
}

// SendWelcomeEmail sends a no-reply welcome message from Lumora to a new user.
func SendWelcomeEmail(cfg config.Config, toEmail, name string) {
	if strings.TrimSpace(name) == "" {
		name = "there"
	}
	_ = SendEmail(cfg, toEmail, "Welcome to Lumora 🦊",
		welcomePlain(name), welcomeHTML(cfg, name))
}

// SendLoginEmail sends a "new sign-in" security alert on login.
func SendLoginEmail(cfg config.Config, toEmail, name string) {
	if strings.TrimSpace(name) == "" {
		name = "there"
	}
	when := time.Now().Format("Mon, 02 Jan 2006 15:04 MST")
	plain := fmt.Sprintf(`Hi %s,

We noticed a new sign-in to your Lumora account on %s.

If this was you, no action is needed. If you don't recognise this activity,
please change your password right away.

— Lumora

This is an automated message. Please do not reply.`, name, when)

	html := simpleHTML(cfg, "New sign-in to your account", []string{
		fmt.Sprintf("Hi %s, we noticed a new sign-in to your Lumora account on <strong>%s</strong>.", name, when),
		"If this was you, you can ignore this email. If you don't recognise it, please change your password right away.",
	}, "Open Lumora", cfg.AppURL)

	_ = SendEmail(cfg, toEmail, "New sign-in to your Lumora account", plain, html)
}

// SendPasswordResetEmail sends a single-use password reset link.
func SendPasswordResetEmail(cfg config.Config, toEmail, name, resetURL string) {
	if strings.TrimSpace(name) == "" {
		name = "there"
	}
	plain := fmt.Sprintf(`Hi %s,

We received a request to reset your Lumora password.

Open this link to choose a new password (it expires in 1 hour):
%s

If you didn't request this, you can safely ignore this email — your password
won't change.

— Lumora

This is an automated message. Please do not reply.`, name, resetURL)

	html := simpleHTML(cfg, "Reset your password", []string{
		fmt.Sprintf("Hi %s, we received a request to reset your Lumora password.", name),
		"Click the button below to choose a new one. This link expires in <strong>1 hour</strong>.",
		"If you didn't request this, you can safely ignore this email — your password won't change.",
	}, "Reset password", resetURL)

	_ = SendEmail(cfg, toEmail, "Reset your Lumora password", plain, html)
}

// SendPaymentEmail sends a receipt after a successful payment. Returns an error
// so callers can confirm delivery (and retry / mark-sent accordingly).
func SendPaymentEmail(cfg config.Config, toEmail, name, itemLabel, amountLabel string) error {
	if strings.TrimSpace(name) == "" {
		name = "there"
	}
	when := time.Now().Format("Mon, 02 Jan 2006 15:04 MST")
	plain := fmt.Sprintf(`Hi %s,

Thank you! We've received your payment.

  Item:   %s
  Amount: %s
  Date:   %s

Your purchase is now active in the app. Enjoy!

— Lumora

This is an automated message. Please do not reply.`, name, itemLabel, amountLabel, when)

	html := simpleHTML(cfg, "Payment received", []string{
		fmt.Sprintf("Hi %s, thank you! We've received your payment.", name),
		fmt.Sprintf("<strong>Item:</strong> %s<br/><strong>Amount:</strong> %s<br/><strong>Date:</strong> %s", itemLabel, amountLabel, when),
		"Your purchase is now active in the app. Enjoy your learning!",
	}, "Back to Lumora", cfg.AppURL)

	return SendEmail(cfg, toEmail, "Your Lumora payment receipt", plain, html)
}

const boundary = "==lumora-mixed-boundary=="

func buildMIME(cfg config.Config, to, subject, plain, html string) string {
	var b strings.Builder
	from := fmt.Sprintf("%s <%s>", cfg.SMTPFromName, cfg.SMTPFrom)

	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: =?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?=\r\n")
	// Signal that this is an automated, no-reply message.
	b.WriteString("Reply-To: " + cfg.SMTPFrom + "\r\n")
	b.WriteString("Auto-Submitted: auto-generated\r\n")
	b.WriteString("X-Auto-Response-Suppress: All\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n\r\n")

	// Plain-text part
	b.WriteString("--" + boundary + "\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	b.WriteString(wrap76(base64.StdEncoding.EncodeToString([]byte(plain))) + "\r\n")

	// HTML part
	b.WriteString("--" + boundary + "\r\n")
	b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	b.WriteString(wrap76(base64.StdEncoding.EncodeToString([]byte(html))) + "\r\n")

	b.WriteString("--" + boundary + "--\r\n")
	return b.String()
}

// wrap76 keeps base64 lines within the SMTP line-length limit.
func wrap76(s string) string {
	var b strings.Builder
	for len(s) > 76 {
		b.WriteString(s[:76] + "\r\n")
		s = s[76:]
	}
	b.WriteString(s)
	return b.String()
}

func welcomePlain(name string) string {
	return fmt.Sprintf(`Hi %s,

Welcome to Lumora! I'm Lumora the fox, and I'll be your guide.

Your account is ready. Here's how to begin:
  1. Pick your language and daily goal
  2. Learn the new words, then practise with quick lessons
  3. Listen, speak and read with your character companions

Open the app and your first lesson is waiting.

— Lumora

This is an automated message. Please do not reply.`, name)
}

func welcomeHTML(cfg config.Config, name string) string {
	// Logo: a hosted image if provided, otherwise an on-brand fox badge.
	logo := `<div style="width:72px;height:72px;line-height:72px;margin:0 auto;border-radius:50%;background:#ffffff;font-size:40px;text-align:center;">🦊</div>`
	if cfg.LogoURL != "" {
		logo = fmt.Sprintf(`<img src="%s" width="72" height="72" alt="Lumora" style="display:block;margin:0 auto;border-radius:50%%;" />`, cfg.LogoURL)
	}

	step := func(n, title, body string) string {
		return fmt.Sprintf(`
      <tr>
        <td style="padding:8px 0;vertical-align:top;width:36px;">
          <div style="width:28px;height:28px;line-height:28px;border-radius:50%%;background:#EDE7F6;color:#6C3FC5;font-weight:800;text-align:center;font-size:14px;">%s</div>
        </td>
        <td style="padding:8px 0;vertical-align:top;">
          <div style="font-weight:700;color:#1A1A2E;font-size:15px;">%s</div>
          <div style="color:#4A4A6A;font-size:14px;">%s</div>
        </td>
      </tr>`, n, title, body)
	}

	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#eceaf3;font-family:-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#eceaf3;padding:24px 0;">
    <tr><td align="center">
      <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:520px;background:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 4px 16px rgba(15,15,36,0.08);">

        <!-- Header -->
        <tr>
          <td style="background:#6C3FC5;padding:32px 24px;text-align:center;">
            %s
            <div style="margin-top:12px;color:#ffffff;font-size:22px;font-weight:800;letter-spacing:-0.5px;">Lumora</div>
            <div style="color:#EDE7F6;font-size:13px;margin-top:2px;">Learn a language. Fall in love with it.</div>
          </td>
        </tr>

        <!-- Body -->
        <tr>
          <td style="padding:28px 28px 8px 28px;">
            <h1 style="margin:0 0 8px 0;color:#1A1A2E;font-size:22px;">Hi %s, welcome aboard!</h1>
            <p style="margin:0;color:#4A4A6A;font-size:15px;line-height:22px;">
              I'm Lumora, your guide. Your account is ready — let's turn a few
              minutes a day into a whole new language.
            </p>
          </td>
        </tr>

        <!-- Steps -->
        <tr>
          <td style="padding:8px 28px 0 28px;">
            <div style="font-size:12px;font-weight:700;letter-spacing:.06em;text-transform:uppercase;color:#9090A0;margin-bottom:4px;">How to start</div>
            <table role="presentation" width="100%%" cellpadding="0" cellspacing="0">
              %s%s%s
            </table>
          </td>
        </tr>

        <!-- CTA -->
        <tr>
          <td style="padding:24px 28px 8px 28px;" align="center">
            <a href="%s" style="display:inline-block;background:#6C3FC5;color:#ffffff;text-decoration:none;font-weight:800;font-size:16px;padding:14px 28px;border-radius:9999px;">
              Start learning
            </a>
          </td>
        </tr>

        <!-- Footer -->
        <tr>
          <td style="padding:20px 28px 28px 28px;text-align:center;">
            <p style="margin:0;color:#9090A0;font-size:12px;line-height:18px;">
              You're receiving this because you created a Lumora account.<br/>
              This is an automated message — please do not reply.
            </p>
          </td>
        </tr>

      </table>
    </td></tr>
  </table>
</body>
</html>`,
		logo, name,
		step("1", "Pick your language &amp; goal", "Choose what to learn and how much time you have."),
		step("2", "Learn, then practise", "Meet the new words first, then lock them in with quick lessons."),
		step("3", "Listen, speak &amp; read", "Train your ear and tongue with your character companions."),
		cfg.AppURL,
	)
}

// simpleHTML is a branded, single-column template used for transactional emails
// (sign-in alerts, receipts). Each entry in `paragraphs` may contain safe HTML.
func simpleHTML(cfg config.Config, heading string, paragraphs []string, ctaText, ctaURL string) string {
	logo := `<div style="width:64px;height:64px;line-height:64px;margin:0 auto;border-radius:50%;background:#ffffff;font-size:34px;text-align:center;">🦊</div>`
	if cfg.LogoURL != "" {
		logo = fmt.Sprintf(`<img src="%s" width="64" height="64" alt="Lumora" style="display:block;margin:0 auto;border-radius:50%%;" />`, cfg.LogoURL)
	}

	var body strings.Builder
	for _, p := range paragraphs {
		body.WriteString(fmt.Sprintf(
			`<p style="margin:0 0 14px 0;color:#4A4A6A;font-size:15px;line-height:22px;">%s</p>`, p))
	}

	cta := ""
	if ctaText != "" && ctaURL != "" {
		cta = fmt.Sprintf(`
        <tr><td style="padding:8px 28px 8px 28px;" align="center">
          <a href="%s" style="display:inline-block;background:#6C3FC5;color:#ffffff;text-decoration:none;font-weight:800;font-size:16px;padding:14px 28px;border-radius:9999px;">%s</a>
        </td></tr>`, ctaURL, ctaText)
	}

	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#eceaf3;font-family:-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#eceaf3;padding:24px 0;">
    <tr><td align="center">
      <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:520px;background:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 4px 16px rgba(15,15,36,0.08);">
        <tr>
          <td style="background:#6C3FC5;padding:28px 24px;text-align:center;">
            %s
            <div style="margin-top:10px;color:#ffffff;font-size:20px;font-weight:800;letter-spacing:-0.5px;">Lumora</div>
          </td>
        </tr>
        <tr>
          <td style="padding:28px 28px 4px 28px;">
            <h1 style="margin:0 0 12px 0;color:#1A1A2E;font-size:20px;">%s</h1>
            %s
          </td>
        </tr>
        %s
        <tr>
          <td style="padding:20px 28px 28px 28px;text-align:center;">
            <p style="margin:0;color:#9090A0;font-size:12px;line-height:18px;">
              This is an automated message from Lumora — please do not reply.
            </p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, logo, heading, body.String(), cta)
}
