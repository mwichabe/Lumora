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

func TestWebOnlyRefusesPhonesAndTablets(t *testing.T) {
	const origin = "https://lumora-learn.netlify.app"
	app := fiber.New()
	app.Get("/exam/paper", WebOnly(origin, origin), func(c *fiber.Ctx) error {
		return c.SendString("paper")
	})

	cases := []struct {
		name, ua, hint string
		want           int
	}{
		{"Windows Chrome", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0 Safari/537.36", "?0", fiber.StatusOK},
		{"Mac Safari", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15", "", fiber.StatusOK},
		{"Linux Firefox", "Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0", "", fiber.StatusOK},
		{"Android Chrome", "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0 Mobile Safari/537.36", "?1", fiber.StatusForbidden},
		{"iPhone Safari", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1", "", fiber.StatusForbidden},
		{"Android tablet", "Mozilla/5.0 (Linux; Android 13; SM-X700) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0 Safari/537.36", "", fiber.StatusForbidden},
		{"mobile client hint only", "Mozilla/5.0", "?1", fiber.StatusForbidden},
	}
	for _, tc := range cases {
		req := httptest.NewRequest("GET", "/exam/paper", nil)
		req.Header.Set("Origin", origin)
		req.Header.Set("User-Agent", tc.ua)
		if tc.hint != "" {
			req.Header.Set("Sec-CH-UA-Mobile", tc.hint)
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
