package controllers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/config"
)

func TestWebAppURL(t *testing.T) {
	a := &AuthController{Cfg: config.Config{
		CORSOrigins: "https://lumora-learn.netlify.app, http://localhost:3000",
		AppURL:      "https://lumora-learn.netlify.app",
	}}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString(a.webAppURL(c)) })

	cases := []struct{ name, origin, want string }{
		{"deployed site", "https://lumora-learn.netlify.app", "https://lumora-learn.netlify.app"},
		{"local dev", "http://localhost:3000", "http://localhost:3000"},
		{"mobile app sends no Origin", "", "https://lumora-learn.netlify.app"},
		{"untrusted origin is ignored", "https://evil.example.com", "https://lumora-learn.netlify.app"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest("GET", "/", nil)
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		got, _ := io.ReadAll(res.Body)
		if string(got) != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}
