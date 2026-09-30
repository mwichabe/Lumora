package utils

import (
	"html"
	"strings"

	"lumora/backend/config"
)

// The shared look of every Lumora email: a night-sky header with the fox, a
// brand stripe, a pill + headline, the message's own rows, and a footer.
//
// Everything is table-based with inline styles — the only layout email clients
// agree on. The <style> block is progressive enhancement (web font, tighter
// padding on phones); the message reads correctly where it's stripped.
//
// The row helpers below take HTML, not text: callers escape anything that came
// from a user (names, labels) with html.EscapeString before passing it in.

// emailLogoURL picks the image shown in the email header: LOGO_URL if set,
// otherwise the app icon the web client already hosts. Returns "" when neither
// is publicly reachable (local dev), and the header falls back to an emoji.
func emailLogoURL(cfg config.Config) string {
	if cfg.LogoURL != "" {
		return cfg.LogoURL
	}
	if strings.HasPrefix(cfg.AppURL, "https://") {
		return strings.TrimRight(cfg.AppURL, "/") + "/icon-192.png"
	}
	return ""
}

// emailDoc describes one message for emailShell.
type emailDoc struct {
	Title     string // <title>; plain text
	Preheader string // inbox preview line next to the subject; plain text
	Pill      string // small label above the headline, e.g. "🔐&nbsp; Password reset"
	Heading   string // the headline (HTML)
	Rows      string // the message body: concatenated <tr> rows from the helpers below
	Footnote  string // why the reader got this email; plain text
}

// emailShell wraps a message's rows in the shared Lumora layout.
func emailShell(cfg config.Config, d emailDoc) string {
	logo := `<div style="width:76px;height:76px;line-height:76px;margin:0 auto;border-radius:22px;background:#6C3FC5;font-size:40px;text-align:center;">🦊</div>`
	if src := emailLogoURL(cfg); src != "" {
		logo = `<img src="` + html.EscapeString(src) + `" width="76" height="76" alt="Lumora" style="display:block;margin:0 auto;border:0;border-radius:22px;" />`
	}
	return strings.NewReplacer(
		"{{logo}}", logo,
		"{{title}}", html.EscapeString(d.Title),
		"{{preheader}}", html.EscapeString(d.Preheader),
		"{{pill}}", d.Pill,
		"{{heading}}", d.Heading,
		"{{rows}}", d.Rows,
		"{{footnote}}", html.EscapeString(d.Footnote),
		"{{appURL}}", html.EscapeString(cfg.AppURL),
	).Replace(emailShellTemplate)
}

// emailParagraph is a block of body copy. align is "center" or "left".
func emailParagraph(content, align string) string {
	return `
        <tr>
          <td class="lm-pad" align="` + align + `" style="padding:14px 44px 0 44px;">
            <p style="margin:0;color:#4A4A6A;font-size:16px;line-height:26px;">` + content + `</p>
          </td>
        </tr>`
}

// emailButton is the primary call to action. It's a padded link inside a
// coloured cell, so it still looks like a button where CSS is stripped.
func emailButton(label, url string) string {
	return `
        <tr>
          <td class="lm-pad" align="center" style="padding:30px 44px 0 44px;">
            <table role="presentation" cellpadding="0" cellspacing="0" border="0" align="center"><tr>
              <td align="center" bgcolor="#6C3FC5" style="background:#6C3FC5;background-image:linear-gradient(135deg,#7B4AD6 0%,#5E33B0 100%);border-radius:9999px;box-shadow:0 8px 20px rgba(108,63,197,0.35);">
                <a href="` + html.EscapeString(url) + `" target="_blank" class="lm-btn" style="display:inline-block;padding:17px 44px;color:#ffffff;font-size:17px;line-height:22px;font-weight:800;text-decoration:none;border-radius:9999px;">` + label + ` &rarr;</a>
              </td>
            </tr></table>
          </td>
        </tr>`
}

// emailChip is a small amber pill for a single fact worth noticing.
func emailChip(content string) string {
	return `
        <tr>
          <td class="lm-pad" align="center" style="padding:18px 44px 0 44px;">
            <div style="display:inline-block;background:#FFF8E7;color:#8A5A00;font-size:13px;line-height:18px;font-weight:700;padding:8px 14px;border-radius:9999px;">` + content + `</div>
          </td>
        </tr>`
}

// emailPanel is a soft lilac box for secondary content (a fallback link, a
// list of steps, a table of details).
func emailPanel(content string) string {
	return `
        <tr>
          <td class="lm-pad" style="padding:32px 44px 0 44px;">
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background:#F7F5FC;border:1px solid #E6DEF7;border-radius:16px;">
              <tr><td style="padding:16px 18px;">` + content + `</td></tr>
            </table>
          </td>
        </tr>`
}

// emailNote is the teal reassurance card: an emoji, a bold line, and detail.
func emailNote(emoji, title, body string) string {
	return `
        <tr>
          <td class="lm-pad" style="padding:14px 44px 0 44px;">
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background:#E0FAF7;border-radius:16px;">
              <tr>
                <td width="46" valign="top" style="padding:16px 0 16px 18px;font-size:22px;line-height:26px;">` + emoji + `</td>
                <td valign="top" style="padding:16px 18px 16px 0;">
                  <div style="color:#0B4F47;font-size:14px;line-height:20px;font-weight:800;">` + title + `</div>
                  <div style="margin-top:2px;color:#1F6F66;font-size:13px;line-height:20px;">` + body + `</div>
                </td>
              </tr>
            </table>
          </td>
        </tr>`
}

// emailDetails renders label/value pairs (sign-in time, receipt lines) as the
// content of an emailPanel. pairs is label, value, label, value, ...
func emailDetails(pairs ...string) string {
	var b strings.Builder
	b.WriteString(`<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">`)
	for i := 0; i+1 < len(pairs); i += 2 {
		border := "border-top:1px solid #E6DEF7;"
		if i == 0 {
			border = ""
		}
		b.WriteString(`<tr>
                  <td style="padding:9px 0;` + border + `color:#9090A0;font-size:12px;line-height:18px;font-weight:800;letter-spacing:.06em;text-transform:uppercase;">` + pairs[i] + `</td>
                  <td align="right" style="padding:9px 0;` + border + `color:#1A1A2E;font-size:14px;line-height:20px;font-weight:800;">` + pairs[i+1] + `</td>
                </tr>`)
	}
	b.WriteString(`</table>`)
	return b.String()
}

// emailSignoff closes the message in Lumora's voice.
func emailSignoff(line string) string {
	return `
        <tr>
          <td class="lm-pad" style="padding:30px 44px 0 44px;">
            <p style="margin:0;color:#4A4A6A;font-size:15px;line-height:24px;">
              ` + line + `<br />
              <strong style="color:#1A1A2E;">Lumora the fox</strong> 🦊
            </p>
          </td>
        </tr>`
}

const emailShellTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <meta name="color-scheme" content="light" />
  <meta name="supported-color-schemes" content="light" />
  <title>{{title}}</title>
  <style>
    @import url('https://fonts.googleapis.com/css2?family=Nunito:wght@400;600;700;800;900&display=swap');
    @media only screen and (max-width:600px) {
      .lm-outer { padding:0 !important; }
      .lm-card { border-radius:0 !important; }
      .lm-pad { padding-left:24px !important; padding-right:24px !important; }
      .lm-h1 { font-size:26px !important; line-height:32px !important; }
      .lm-btn { display:block !important; }
    }
  </style>
</head>
<body style="margin:0;padding:0;background:#eceaf3;-webkit-text-size-adjust:100%;">
  <!-- Preheader: the grey preview line next to the subject in the inbox. -->
  <div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;font-size:1px;line-height:1px;">
    {{preheader}}
    &#847;&zwnj;&nbsp;&#847;&zwnj;&nbsp;&#847;&zwnj;&nbsp;&#847;&zwnj;&nbsp;&#847;&zwnj;&nbsp;&#847;&zwnj;&nbsp;&#847;&zwnj;&nbsp;&#847;&zwnj;&nbsp;
  </div>

  <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" class="lm-outer" style="background:#eceaf3;padding:32px 12px;">
    <tr><td align="center">
      <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" class="lm-card" style="max-width:560px;background:#ffffff;border-radius:24px;overflow:hidden;box-shadow:0 12px 40px rgba(58,31,138,0.14);font-family:'Nunito',-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;">

        <!-- Header: night-sky gradient, solid colour where gradients are stripped -->
        <tr>
          <td align="center" bgcolor="#2A1669" style="background:#2A1669;background-image:linear-gradient(145deg,#0F0F24 0%,#3A1F8A 52%,#6C3FC5 100%);padding:40px 24px 36px 24px;">
            {{logo}}
            <div style="margin-top:14px;color:#ffffff;font-size:24px;line-height:28px;font-weight:900;letter-spacing:-0.5px;">Lumora</div>
            <div style="margin-top:4px;color:#C9B8F2;font-size:13px;line-height:18px;font-weight:600;">Learn a language. Fall in love with it.</div>
          </td>
        </tr>
        <!-- Brand stripe -->
        <tr>
          <td style="padding:0;">
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0"><tr>
              <td width="34%" height="5" bgcolor="#F5A623" style="background:#F5A623;font-size:0;line-height:0;">&nbsp;</td>
              <td width="33%" height="5" bgcolor="#00C2A8" style="background:#00C2A8;font-size:0;line-height:0;">&nbsp;</td>
              <td width="33%" height="5" bgcolor="#FF5C5C" style="background:#FF5C5C;font-size:0;line-height:0;">&nbsp;</td>
            </tr></table>
          </td>
        </tr>

        <!-- Headline -->
        <tr>
          <td class="lm-pad" align="center" style="padding:40px 44px 0 44px;">
            <div style="display:inline-block;background:#EDE7F6;color:#6C3FC5;font-size:12px;line-height:16px;font-weight:800;letter-spacing:.08em;text-transform:uppercase;padding:7px 14px;border-radius:9999px;">{{pill}}</div>
            <h1 class="lm-h1" style="margin:18px 0 0 0;color:#1A1A2E;font-size:30px;line-height:36px;font-weight:900;letter-spacing:-0.6px;">{{heading}}</h1>
          </td>
        </tr>
{{rows}}

        <!-- Footer -->
        <tr>
          <td class="lm-pad" style="padding:30px 44px 36px 44px;">
            <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0"><tr>
              <td style="border-top:1px solid #EBEBEB;padding-top:20px;" align="center">
                <p style="margin:0;color:#9090A0;font-size:12px;line-height:19px;">
                  {{footnote}}<br />
                  This is an automated message — please don't reply.
                </p>
                <p style="margin:10px 0 0 0;font-size:12px;line-height:19px;">
                  <a href="{{appURL}}" target="_blank" style="color:#6C3FC5;font-weight:800;text-decoration:none;">Open Lumora</a>
                </p>
              </td>
            </tr></table>
          </td>
        </tr>

      </table>
    </td></tr>
  </table>
</body>
</html>`
