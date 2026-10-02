package utils

import (
	"regexp"
	"strings"
)

// emailPattern is the address format accepted everywhere — the same pattern
// is used by the web (frontend/lib/validation.ts) and mobile
// (mobile/lib/core/validation.dart) clients, so all three agree.
//
//   - local part: letters, digits and !#$%&'*+/=?^_`{|}~- in dot-separated
//     runs (no leading, trailing or doubled dots)
//   - domain: one or more labels of letters, digits and inner hyphens (a label
//     can't start or end with a hyphen), then a TLD of 2–63 letters
var emailPattern = regexp.MustCompile(
	"^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*" +
		"@([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z]{2,63}$")

// NormaliseEmail trims and lower-cases an address, the form it's stored in.
func NormaliseEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// EmailError returns a message saying what's wrong with an address, or "" if
// it's valid. Pass it a normalised address.
func EmailError(email string) string {
	switch {
	case email == "":
		return "Enter your email address."
	case len(email) > 254:
		return "That email address is too long."
	case strings.Count(email, "@") != 1:
		return "Enter a valid email address, like name@example.com."
	}
	local := email[:strings.Index(email, "@")]
	if len(local) > 64 {
		return "The part before the @ is too long."
	}
	if !emailPattern.MatchString(email) {
		return "Enter a valid email address, like name@example.com."
	}
	return ""
}
