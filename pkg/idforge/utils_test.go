package idforge

import (
	"testing"
)

func TestGenerateSecureToken(t *testing.T) {
	length := 16
	token, err := GenerateSecureToken(length)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(token) != length {
		t.Errorf("Expected token length %d, got %d", length, len(token))
	}
}

func TestMustGenerateSecureToken(t *testing.T) {
	length := 32
	token := MustGenerateSecureToken(length)
	if len(token) != length {
		t.Errorf("Expected token length %d, got %d", length, len(token))
	}
}

func TestMustGenerateSecureTokenPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but no panic occurred")
		}
	}()
	MustGenerateSecureToken(-1)
}

func TestIsValidID(t *testing.T) {
	testCases := []struct {
		id       string
		alphabet string
		size     int
		expected bool
	}{
		{"ABC123", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6, true},
		{"ABC123", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 5, false},
		{"ABC123", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", 6, false},
		{"abc123", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6, false},
	}

	for _, tc := range testCases {
		result := IsValidID(tc.id, tc.alphabet, tc.size)
		if result != tc.expected {
			t.Errorf("Expected %v for ID '%s' with alphabet '%s' and size %d, got %v",
				tc.expected, tc.id, tc.alphabet, tc.size, result)
		}
	}
}
