package middleware

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AllowedOrigin returns the request's Origin when it is one of the web app's
// origins (the comma-separated CORS allow-list), and "" otherwise. Browsers
// attach Origin to every cross-origin call and page scripts can't forge it;
// native apps don't send one at all.
func AllowedOrigin(c *fiber.Ctx, corsOrigins string) string {
	origin := strings.TrimRight(c.Get(fiber.HeaderOrigin), "/")
	if origin == "" {
		return ""
	}
	for _, allowed := range strings.Split(corsOrigins, ",") {
		if strings.EqualFold(origin, strings.TrimRight(strings.TrimSpace(allowed), "/")) {
			return origin
		}
	}
	return ""
}

// WebOnly rejects requests that don't come from the web app. It guards the
// exam: sitting one requires camera + screen-share proctoring that only the
// browser client implements, so the mobile app may pay for an attempt but not
// take it. This is a product rule, not a security boundary — a determined
// caller can set the header by hand — but it holds for every shipped client,
// including older app builds that still contain the exam screens.
//
// The web app must also be open on a laptop or desktop: a phone or tablet
// browser is refused the same way (see IsHandheld). The page checks this
// itself; this is the server's half of the rule.
func WebOnly(corsOrigins, appURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if AllowedOrigin(c, corsOrigins) == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Exams can only be taken on the web. Sign in at " + appURL + " on a computer to take yours.",
				"webOnly": true,
			})
		}
		if IsHandheld(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":       "Exams can only be taken on a laptop or desktop computer, not a phone or tablet. Open " + appURL + " on a computer to take yours.",
				"desktopOnly": true,
			})
		}
		return c.Next()
	}
}

// handheldUA matches the user agents of phone and tablet browsers. iPadOS
// Safari reports itself as a Mac, so it slips through here — the web page
// catches it by its touch screen.
var handheldUA = regexp.MustCompile(`(?i)android|iphone|ipad|ipod|mobile|tablet|silk|kindle|playbook|bb10|opera mini|iemobile|windows phone`)

// IsHandheld reports whether the request comes from a phone or tablet browser:
// the Sec-CH-UA-Mobile client hint when the browser sends it, else the user
// agent string.
func IsHandheld(c *fiber.Ctx) bool {
	if c.Get("Sec-CH-UA-Mobile") == "?1" {
		return true
	}
	return handheldUA.MatchString(c.Get(fiber.HeaderUserAgent))
}
