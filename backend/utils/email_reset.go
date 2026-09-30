package utils

import (
	"fmt"
	"html"

	"lumora/backend/config"
)

// The password-reset email — the one message a locked-out user is actively
// waiting for. The button has to be unmissable, the link has to survive
// clients that strip buttons, and the "wasn't you?" reassurance has to be
// impossible to overlook.

func passwordResetPlain(name, resetURL string) string {
	return fmt.Sprintf(`Hi %s,

We received a request to reset the password for your Lumora account.

Choose a new password here (the link expires in 1 hour and works once):
%s

Didn't ask for this? You can safely ignore this email — your password
won't change, and nobody can get into your account with this message alone.

— Lumora the fox

This is an automated message. Please do not reply.`, name, resetURL)
}

func passwordResetHTML(cfg config.Config, name, resetURL string) string {
	// name comes from the user's profile and resetURL is placed in attributes,
	// so both are escaped before they touch the markup.
	safeName, safeURL := html.EscapeString(name), html.EscapeString(resetURL)

	fallback := `<div style="color:#1A1A2E;font-size:13px;line-height:18px;font-weight:800;">Button not working?</div>
                <div style="margin-top:2px;color:#4A4A6A;font-size:13px;line-height:19px;">Copy this link and paste it into your browser:</div>
                <div style="margin-top:8px;font-family:'SFMono-Regular',Consolas,'Liberation Mono',Menlo,monospace;font-size:12px;line-height:18px;word-break:break-all;">
                  <a href="` + safeURL + `" target="_blank" style="color:#6C3FC5;text-decoration:underline;">` + safeURL + `</a>
                </div>`

	return emailShell(cfg, emailDoc{
		Title:     "Reset your Lumora password",
		Preheader: "Your reset link is inside — it's valid for 1 hour. Didn't ask for it? Just ignore this.",
		Pill:      "🔐&nbsp; Password reset",
		Heading:   "Let's get you back in",
		Rows: emailParagraph(`Hi <strong style="color:#1A1A2E;">`+safeName+`</strong>, we got a request to reset the password
              for your Lumora account. Tap the button below to choose a new one.`, "center") +
			emailButton("Reset my password", resetURL) +
			emailChip("⏳&nbsp; Expires in 1 hour &middot; works only once") +
			emailPanel(fallback) +
			emailNote("🛡️", "Didn't ask for this?",
				"No action needed — just ignore this email. Your password stays exactly as it is, and nobody can get into your account with this message alone.") +
			emailSignoff("See you back in your lessons,"),
		Footnote: "You're receiving this because a password reset was requested for your Lumora account.",
	})
}
