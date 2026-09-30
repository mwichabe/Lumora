package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestWebOnly(t *testing.T) {
	const origins = "https://lumora-learn.netlify.app, http://localhost:3000"
	app := fiber.New()
	app.Get("/exam/paper", WebOnly(origins, "https://lumora-learn.netlify.app"), func(c *fiber.Ctx) error {
		return c.SendString("paper")
	})

	cases := []struct {
		name, origin string
		want         int
	}{
		{"web app", "https://lumora-learn.netlify.app", fiber.StatusOK},
		{"local web dev", "http://localhost:3000", fiber.StatusOK},
		{"mobile app sends no Origin", "", fiber.StatusForbidden},
		{"another site", "https://evil.example.com", fiber.StatusForbidden},
	}
	for _, tc := range cases {
		req := httptest.NewRequest("GET", "/exam/paper", nil)
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if res.StatusCode != tc.want {
			t.Errorf("%s: status %d, want %d", tc.name, res.StatusCode, tc.want)
		}
	}
}
