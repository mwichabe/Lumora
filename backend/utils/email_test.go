package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lumora/backend/config"
)

// stubResend points the sender at a local server for the duration of a test.
func stubResend(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	prev := resendEndpoint
	resendEndpoint = srv.URL
	t.Cleanup(func() {
		resendEndpoint = prev
		srv.Close()
	})
}

func TestSendEmailUsesResendWhenKeySet(t *testing.T) {
	var gotAuth string
	var got resendPayload
	stubResend(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("payload is not valid JSON: %v", err)
		}
		w.Write([]byte(`{"id":"abc"}`))
	})

	// SMTP is deliberately left unconfigured: Resend must not depend on it.
	cfg := config.Config{ResendAPIKey: "re_test", ResendFrom: "Lumora <no-reply@example.com>"}
	if err := SendEmail(cfg, "learner@example.com", "Hello", "plain body", "<p>html body</p>"); err != nil {
		t.Fatalf("SendEmail returned an error: %v", err)
	}

	if gotAuth != "Bearer re_test" {
		t.Errorf("Authorization = %q, want the bearer key", gotAuth)
	}
	if got.From != cfg.ResendFrom || len(got.To) != 1 || got.To[0] != "learner@example.com" {
		t.Errorf("unexpected envelope: from=%q to=%v", got.From, got.To)
	}
	if got.Subject != "Hello" || got.Text != "plain body" || got.HTML != "<p>html body</p>" {
		t.Errorf("unexpected content: %+v", got)
	}
}

func TestSendEmailSurfacesResendRejection(t *testing.T) {
	stubResend(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"The example.com domain is not verified."}`))
	})

	cfg := config.Config{ResendAPIKey: "re_test", ResendFrom: "Lumora <no-reply@example.com>"}
	err := SendEmail(cfg, "learner@example.com", "Hello", "plain", "<p>html</p>")
	if err == nil {
		t.Fatal("expected an error for a 403 reply")
	}
	if !strings.Contains(err.Error(), "403") || !strings.Contains(err.Error(), "not verified") {
		t.Errorf("error should carry Resend's status and reason, got: %v", err)
	}
}

func TestPasswordResetHTML(t *testing.T) {
	cfg := config.Config{AppURL: "https://lumora-learn.netlify.app"}
	resetURL := "https://lumora-learn.netlify.app/reset-password?token=abc123"
	out := passwordResetHTML(cfg, `<script>alert(1)</script>`, resetURL)

	if strings.Contains(out, "<script>") {
		t.Error("the user's name must be HTML-escaped")
	}
	// Once on the button, once as the link target of the fallback, once as its text.
	if n := strings.Count(out, resetURL); n != 3 {
		t.Errorf("reset URL appears %d times, want 3", n)
	}
	if strings.Contains(out, "{{") {
		t.Error("template has an unfilled placeholder")
	}
	if !strings.Contains(out, "https://lumora-learn.netlify.app/icon-192.png") {
		t.Error("expected the hosted app icon as the header logo")
	}
}

func TestEmailLogoURL(t *testing.T) {
	cases := []struct {
		name string
		cfg  config.Config
		want string
	}{
		{"explicit logo wins", config.Config{LogoURL: "https://cdn.example.com/fox.png", AppURL: "https://app.example.com"}, "https://cdn.example.com/fox.png"},
		{"hosted app icon", config.Config{AppURL: "https://app.example.com/"}, "https://app.example.com/icon-192.png"},
		{"localhost is unreachable from an inbox", config.Config{AppURL: "http://localhost:3000"}, ""},
	}
	for _, c := range cases {
		if got := emailLogoURL(c.cfg); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
