package middleware

import (
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
func WebOnly(corsOrigins, appURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if AllowedOrigin(c, corsOrigins) == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Exams can only be taken on the web. Sign in at " + appURL + " on a computer to take yours.",
				"webOnly": true,
			})
		}
		return c.Next()
	}
}
