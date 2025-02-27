package idforge

import (
	"testing"
)

func TestNewIDValidator(t *testing.T) {
	validator := NewIDValidator()
	if validator == nil {
		t.Error("Expected a non-nil validator")
	}
	if validator.minLength != 8 {
		t.Errorf("Expected minLength to be 8, got %d", validator.minLength)
	}
	if validator.maxLength != 128 {
		t.Errorf("Expected maxLength to be 128, got %d", validator.maxLength)
	}
}

func TestWithMinLength(t *testing.T) {
	validator := NewIDValidator(WithMinLength(12))
	if validator.minLength != 12 {
		t.Errorf("Expected minLength to be 12, got %d", validator.minLength)
	}
}

func TestWithMaxLength(t *testing.T) {
	validator := NewIDValidator(WithMaxLength(64))
	if validator.maxLength != 64 {
		t.Errorf("Expected maxLength to be 64, got %d", validator.maxLength)
	}
}

// func TestValidate(t *testing.T) {
// 	validator := NewIDValidator()

// 	testCases := []struct {
// 		id       string
// 		expected error
// 	}{
// 		{"abcABC123", nil},
// 		{"abc", ErrIDTooShort},
// 		{"abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789", ErrIDTooLong},
// 		{"abcdef", ErrWeakID},
// 		{"ABCDEF", ErrWeakID},
// 		{"123456", ErrWeakID},
// 	}

// 	for _, tc := range testCases {
// 		err := validator.Validate(tc.id)
// 		if err != tc.expected {
// 			t.Errorf("Expected error '%v' for ID '%s', got '%v'", tc.expected, tc.id, err)
// 		}
// 	}
// }

func TestAddForbiddenPattern(t *testing.T) {
	validator := NewIDValidator()
	AddForbiddenPattern(`\d{4}`)(validator)

	err := validator.Validate("abcABC1234")
	if err == nil {
		t.Error("Expected an error for forbidden pattern, got nil")
	}
}

func TestSecureCompare(t *testing.T) {
	id1 := "abcABC123"
	id2 := "abcABC123"
	id3 := "xyzXYZ789"

	if !SecureCompare(id1, id2) {
		t.Errorf("Expected '%s' and '%s' to be equal", id1, id2)
	}
	if SecureCompare(id1, id3) {
		t.Errorf("Expected '%s' and '%s' to be different", id1, id3)
	}
}

func TestSanitizeID(t *testing.T) {
	id := "abc123!@#"
	expected := "abc123"

	sanitized := SanitizeID(id, "abcdefghijklmnopqrstuvwxyz0123456789")
	if sanitized != expected {
		t.Errorf("Expected sanitized ID to be '%s', got '%s'", expected, sanitized)
	}
}

func TestGenerateComplexityReport(t *testing.T) {
	id := "abcABC123!@#"
	report := GenerateComplexityReport(id)

	expectedLength := 12
	if report["length"].(int) != expectedLength {
		t.Errorf("Expected length to be %d, got %d", expectedLength, report["length"].(int))
	}

	expectedLowercase := 3
	if report["complexity"].(map[string]int)["lowercase"] != expectedLowercase {
		t.Errorf("Expected lowercase count to be %d, got %d",
			expectedLowercase, report["complexity"].(map[string]int)["lowercase"])
	}

	expectedUppercase := 3
	if report["complexity"].(map[string]int)["uppercase"] != expectedUppercase {
		t.Errorf("Expected uppercase count to be %d, got %d",
			expectedUppercase, report["complexity"].(map[string]int)["uppercase"])
	}

	expectedDigits := 3
	if report["complexity"].(map[string]int)["digits"] != expectedDigits {
		t.Errorf("Expected digits count to be %d, got %d",
			expectedDigits, report["complexity"].(map[string]int)["digits"])
	}

	expectedSymbols := 3
	if report["complexity"].(map[string]int)["symbols"] != expectedSymbols {
		t.Errorf("Expected symbols count to be %d, got %d",
			expectedSymbols, report["complexity"].(map[string]int)["symbols"])
	}
}
