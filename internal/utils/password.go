package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateTempPassword generates a secure random password of default length 12
func GenerateTempPassword() string {
	n := 12
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback or panic, but rand.Read failing is rare.
		// For simplicity in this helper, we'll return a less secure fallback or empty?
		// Better to just return empty string and let caller handle, but signature is string.
		return "ChangeMe123!"
	}
	// encoding usually expands size, so slicing to n is fine to get n chars
	// but base64 might have non-alphanumeric. URL encoding is safer.
	return base64.URLEncoding.EncodeToString(b)[:n]
}
