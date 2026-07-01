package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomToken returns a cryptographically-random hex string of `nBytes` bytes
// (so the string is 2*nBytes chars). Used for password-reset links.
func RandomToken(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		// Extremely unlikely; fall back to an empty token which the caller's
		// unique constraint / expiry will still reject if reused.
		return ""
	}
	return hex.EncodeToString(b)
}
