package utils

import (
	"encoding/base64"
	"fmt"
	"html"
	"log"
	"net/smtp"
	"strings"
	"time"

	"lumora/backend/config"
)

// SendEmail is the core sender. It delivers through the Resend HTTP API when
// RESEND_API_KEY is set, and otherwise builds a MIME multipart message and
// sends it via Gmail SMTP (STARTTLS on port 587) using the App Password from
// config. Safe to call in a goroutine — failures are logged, never fatal, and
// if neither transport is configured it simply no-ops.
func SendEmail(cfg config.Config, toEmail, subject, plain, html string) error {
	if cfg.ResendAPIKey != "" {
		if err := sendViaResend(cfg, toEmail, subject, plain, html); err != nil {
			log.Printf("[email] failed to send '%s' to %s: %v", subject, toEmail, err)
			return err
		}
		log.Printf("[email] sent '%s' to %s", subject, toEmail)
		return nil
	}
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		log.Printf("[email] no email transport configured — skipping '%s' to %s", subject, toEmail)
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

// SendLoginEmail welcomes a user back on sign-in. It doubles as the security
// alert: it says when the sign-in happened and how to lock the account if it
// wasn't them.
func SendLoginEmail(cfg config.Config, toEmail, name string) {
	if strings.TrimSpace(name) == "" {
		name = "there"
	}
	when := time.Now().Format("Mon, 02 Jan 2006 15:04 MST")
	_ = SendEmail(cfg, toEmail, "Welcome back to Lumora 👋",
		loginPlain(cfg, name, when), loginHTML(cfg, toEmail, name, when))
}

// SendPasswordResetEmail sends a single-use password reset link.
func SendPasswordResetEmail(cfg config.Config, toEmail, name, resetURL string) {
	if strings.TrimSpace(name) == "" {
		name = "there"
	}
	_ = SendEmail(cfg, toEmail, "Reset your Lumora password",
		passwordResetPlain(name, resetURL), passwordResetHTML(cfg, name, resetURL))
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

	return SendEmail(cfg, toEmail, "Your Lumora payment receipt", plain,
		paymentHTML(cfg, name, itemLabel, amountLabel, when))
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
	step := func(n, title, body string) string {
		return `<tr>
                  <td width="40" valign="top" style="padding:9px 0;">
                    <div style="width:28px;height:28px;line-height:28px;border-radius:50%;background:#6C3FC5;color:#ffffff;font-weight:800;text-align:center;font-size:14px;">` + n + `</div>
                  </td>
                  <td valign="top" style="padding:9px 0;">
                    <div style="color:#1A1A2E;font-size:15px;line-height:20px;font-weight:800;">` + title + `</div>
                    <div style="margin-top:2px;color:#4A4A6A;font-size:13px;line-height:20px;">` + body + `</div>
                  </td>
                </tr>`
	}
	steps := `<div style="color:#9090A0;font-size:12px;line-height:18px;font-weight:800;letter-spacing:.06em;text-transform:uppercase;">How to start</div>
                <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="margin-top:4px;">` +
		step("1", "Pick your language &amp; goal", "Choose what to learn and how much time you have.") +
		step("2", "Learn, then practise", "Meet the new words first, then lock them in with quick lessons.") +
		step("3", "Listen, speak &amp; read", "Train your ear and tongue with your character companions.") +
		`</table>`

	return emailShell(cfg, emailDoc{
		Title:     "Welcome to Lumora",
		Preheader: "Your account is ready — your first lesson is waiting.",
		Pill:      "🎉&nbsp; Welcome aboard",
		Heading:   "Hi " + html.EscapeString(name) + ", welcome to Lumora!",
		Rows: emailParagraph(`I'm Lumora, your guide. Your account is ready — let's turn a few
              minutes a day into a whole new language.`, "center") +
			emailPanel(steps) +
			emailButton("Start learning", cfg.AppURL) +
			emailNote("🔥", "Start your streak today",
				"A few minutes every day beats an hour once a week. Finish one lesson today and your streak begins.") +
			emailSignoff("See you in your first lesson,"),
		Footnote: "You're receiving this because you created a Lumora account.",
	})
}

func loginPlain(cfg config.Config, name, when string) string {
	return fmt.Sprintf(`Welcome back, %s!

You just signed in to your Lumora account on %s. Your lessons are right
where you left them:
%s

Wasn't you? Reset your password right away:
%s/forgot-password

— Lumora the fox

This is an automated message. Please do not reply.`, name, when, cfg.AppURL, cfg.AppURL)
}

func loginHTML(cfg config.Config, toEmail, name, when string) string {
	resetLink := `<a href="` + html.EscapeString(cfg.AppURL+"/forgot-password") + `" target="_blank" style="color:#0B4F47;font-weight:800;">reset your password</a>`

	return emailShell(cfg, emailDoc{
		Title:     "Welcome back to Lumora",
		Preheader: "You just signed in. Your lessons are right where you left them.",
		Pill:      "👋&nbsp; Welcome back",
		Heading:   "Good to see you, " + html.EscapeString(name) + "!",
		Rows: emailParagraph(`You just signed in to your Lumora account. Your lessons are right
              where you left them — a quick one today keeps your streak alive.`, "center") +
			emailButton("Continue learning", cfg.AppURL) +
			emailPanel(emailDetails(
				"Signed in", html.EscapeString(when),
				"Account", html.EscapeString(toEmail),
			)) +
			emailNote("🛡️", "Wasn't you?",
				"If you don't recognise this sign-in, "+resetLink+" right away to lock everyone else out.") +
			emailSignoff("Happy learning,"),
		Footnote: "You're receiving this because someone signed in to your Lumora account.",
	})
}

func paymentHTML(cfg config.Config, name, itemLabel, amountLabel, when string) string {
	return emailShell(cfg, emailDoc{
		Title:     "Your Lumora payment receipt",
		Preheader: "We've received your payment — your purchase is active.",
		Pill:      "✅&nbsp; Payment received",
		Heading:   "Thank you, " + html.EscapeString(name) + "!",
		Rows: emailParagraph(`We've received your payment and your purchase is now active in the app.`, "center") +
			emailPanel(emailDetails(
				"Item", html.EscapeString(itemLabel),
				"Amount", html.EscapeString(amountLabel),
				"Date", html.EscapeString(when),
			)) +
			emailButton("Back to Lumora", cfg.AppURL) +
			emailSignoff("Enjoy your learning,"),
		Footnote: "You're receiving this because a payment was made on your Lumora account.",
	})
}
