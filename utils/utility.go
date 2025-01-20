package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomString generates a cryptographically secure random string of the specified length
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length: must be greater than 0")
	}

	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Encode to a base64 string to ensure safe and readable output
	return base64.RawURLEncoding.EncodeToString(randomBytes)[:length], nil
}
