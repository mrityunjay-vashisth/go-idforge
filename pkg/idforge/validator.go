package idforge

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Predefined error types for more specific error handling
	ErrIDTooShort       = errors.New("ID is too short")
	ErrIDTooLong        = errors.New("ID is too long")
	ErrInvalidCharacter = errors.New("ID contains invalid characters")
	ErrWeakID           = errors.New("generated ID does not meet complexity requirements")
)

// IDValidator provides advanced validation capabilities
type IDValidator struct {
	minLength         int
	maxLength         int
	requiredCharSet   []CharacterSetRequirement
	forbiddenPatterns []*regexp.Regexp
}

// CharacterSetRequirement defines rules for character set composition
type CharacterSetRequirement struct {
	CharSet     string
	MinCount    int
	Description string
}

// NewIDValidator creates a configurable ID validator
func NewIDValidator(opts ...func(*IDValidator)) *IDValidator {
	validator := &IDValidator{
		minLength: 8,   // Sensible default minimum
		maxLength: 128, // Reasonable maximum
		requiredCharSet: []CharacterSetRequirement{
			{
				CharSet:     "0123456789",
				MinCount:    1,
				Description: "at least one digit",
			},
			{
				CharSet:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
				MinCount:    1,
				Description: "at least one uppercase letter",
			},
			{
				CharSet:     "abcdefghijklmnopqrstuvwxyz",
				MinCount:    1,
				Description: "at least one lowercase letter",
			},
		},
		forbiddenPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)password`),
			regexp.MustCompile(`^0+$`),
			regexp.MustCompile(`^9+$`),
		},
	}

	// Apply custom options
	for _, opt := range opts {
		opt(validator)
	}

	return validator
}

// WithMinLength sets minimum ID length
func WithMinLength(min int) func(*IDValidator) {
	return func(v *IDValidator) {
		if min > 0 {
			v.minLength = min
		}
	}
}

// WithMaxLength sets maximum ID length
func WithMaxLength(max int) func(*IDValidator) {
	return func(v *IDValidator) {
		if max > 0 {
			v.maxLength = max
		}
	}
}

// AddForbiddenPattern allows adding custom regex patterns to reject
func AddForbiddenPattern(pattern string) func(*IDValidator) {
	return func(v *IDValidator) {
		regex, err := regexp.Compile(pattern)
		if err == nil {
			v.forbiddenPatterns = append(v.forbiddenPatterns, regex)
		}
	}
}

// Validate performs comprehensive ID validation
func (v *IDValidator) Validate(id string) error {
	// Length check
	if len(id) < v.minLength {
		return fmt.Errorf("%w: expected at least %d characters", ErrIDTooShort, v.minLength)
	}
	if len(id) > v.maxLength {
		return fmt.Errorf("%w: expected at most %d characters", ErrIDTooLong, v.maxLength)
	}

	// Forbidden pattern check
	for _, pattern := range v.forbiddenPatterns {
		if pattern.MatchString(id) {
			return fmt.Errorf("ID matches forbidden pattern: %v", pattern)
		}
	}

	// Character set requirements
	charCounts := make(map[rune]int)
	for _, char := range id {
		charCounts[char]++
	}

	var missingRequirements []string
	for _, req := range v.requiredCharSet {
		count := 0
		for _, char := range req.CharSet {
			count += charCounts[char]
		}
		if count < req.MinCount {
			missingRequirements = append(missingRequirements, req.Description)
		}
	}

	if len(missingRequirements) > 0 {
		return fmt.Errorf("%w: ID must contain %s",
			ErrWeakID,
			strings.Join(missingRequirements, ", "))
	}

	return nil
}

// SecureCompare provides timing-safe comparison of IDs
func SecureCompare(a, b string) bool {
	// Constant-time comparison to prevent timing attacks
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// SanitizeID cleans and normalizes an ID
func SanitizeID(id string, allowedChars string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(allowedChars, r) {
			return r
		}
		return -1
	}, id)
}

// GenerateComplexityReport provides insights into ID complexity
func GenerateComplexityReport(id string) map[string]interface{} {
	report := map[string]interface{}{
		"length": len(id),
		"complexity": map[string]int{
			"lowercase": 0,
			"uppercase": 0,
			"digits":    0,
			"symbols":   0,
		},
		"entropy": calculateEntropy(id),
	}

	for _, char := range id {
		switch {
		case unicode.IsLower(char):
			report["complexity"].(map[string]int)["lowercase"]++
		case unicode.IsUpper(char):
			report["complexity"].(map[string]int)["uppercase"]++
		case unicode.IsDigit(char):
			report["complexity"].(map[string]int)["digits"]++
		default:
			report["complexity"].(map[string]int)["symbols"]++
		}
	}

	return report
}

// calculateEntropy estimates the Shannon entropy of the ID
func calculateEntropy(id string) float64 {
	charCount := make(map[rune]int)
	for _, char := range id {
		charCount[char]++
	}

	entropy := 0.0
	idLength := len(id)
	for _, count := range charCount {
		prob := float64(count) / float64(idLength)
		entropy -= prob * math.Log2(prob)
	}

	return entropy
}
