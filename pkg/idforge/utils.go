package idforge

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

// GenerateSecureToken creates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(b)[:length], nil
}

// MustGenerateSecureToken generates a token, panicking on error
func MustGenerateSecureToken(length int) string {
	token, err := GenerateSecureToken(length)
	if err != nil {
		panic(err)
	}
	return token
}

// IsValidID checks if the ID follows standard generation rules
func IsValidID(id string, alphabet string, size int) bool {
	if len(id) != size {
		return false
	}

	for _, char := range id {
		if !strings.ContainsRune(alphabet, char) {
			return false
		}
	}

	return true
}
