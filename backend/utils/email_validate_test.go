package utils

import "testing"

func TestEmailValidation(t *testing.T) {
	valid := []string{
		"name@example.com",
		"first.last@example.co.uk",
		"user+tag@sub.domain.io",
		"o'brien@example.ie",
		"a@b.co",
		"x_y-z@my-domain.com",
	}
	invalid := []string{
		"",
		"plainaddress",
		"@example.com",
		"name@",
		"name@example",
		"name@example.c",
		"name@@example.com",
		"name@ex@ample.com",
		".name@example.com",
		"name.@example.com",
		"na..me@example.com",
		"name@-example.com",
		"name@example-.com",
		"name@example..com",
		"name@.example.com",
		"name @example.com",
		"name@exa mple.com",
		"name@example.123",
		"name@example.com.",
	}
	for _, e := range valid {
		if msg := EmailError(NormaliseEmail(e)); msg != "" {
			t.Errorf("%q rejected: %s", e, msg)
		}
	}
	for _, e := range invalid {
		if EmailError(NormaliseEmail(e)) == "" {
			t.Errorf("%q accepted", e)
		}
	}
}
