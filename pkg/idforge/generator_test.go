package idforge

import (
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	// Test default generator
	defaultGen := New()
	if defaultGen.alphabet != DefaultAlphabet {
		t.Errorf("Default alphabet incorrect. Expected %s, got %s",
			DefaultAlphabet, defaultGen.alphabet)
	}
	if defaultGen.size != DefaultSize {
		t.Errorf("Default size incorrect. Expected %d, got %d",
			DefaultSize, defaultGen.size)
	}

	// Test generator with custom options
	customGen := New(
		WithAlphabet("ABC"),
		WithSize(10),
	)
	if customGen.alphabet != "ABC" {
		t.Errorf("Custom alphabet not set correctly. Got %s", customGen.alphabet)
	}
	if customGen.size != 10 {
		t.Errorf("Custom size not set correctly. Got %d", customGen.size)
	}
}

func TestGenerate(t *testing.T) {
	// Test default generator
	gen := New()

	// Generate multiple IDs to ensure uniqueness and basic properties
	generatedIDs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := gen.Generate()
		if err != nil {
			t.Fatalf("Unexpected error generating ID: %v", err)
		}

		// Check ID length
		if len(id) != DefaultSize {
			t.Errorf("Generated ID length incorrect. Expected %d, got %d",
				DefaultSize, len(id))
		}

		// Check for characters from alphabet
		for _, char := range id {
			if !strings.ContainsRune(DefaultAlphabet, char) {
				t.Errorf("ID contains character not in alphabet: %c", char)
			}
		}

		// Check for uniqueness
		if generatedIDs[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		generatedIDs[id] = true
	}
}

func TestGenerateWithCustomOptions(t *testing.T) {
	// Test generation with custom alphabet and size
	customAlphabet := "XYZ"
	customSize := 8
	gen := New(
		WithAlphabet(customAlphabet),
		WithSize(customSize),
	)

	for i := 0; i < 50; i++ {
		id, err := gen.Generate()
		if err != nil {
			t.Fatalf("Unexpected error generating ID: %v", err)
		}

		// Check ID length
		if len(id) != customSize {
			t.Errorf("Generated ID length incorrect. Expected %d, got %d",
				customSize, len(id))
		}

		// Check for characters from custom alphabet
		for _, char := range id {
			if !strings.ContainsRune(customAlphabet, char) {
				t.Errorf("ID contains character not in custom alphabet: %c", char)
			}
		}
	}
}

func TestValidate(t *testing.T) {
	// Test default generator validation
	gen := New()

	// Test valid IDs
	validID, _ := gen.Generate()
	if !gen.Validate(validID) {
		t.Errorf("Valid ID %s failed validation", validID)
	}

	// Test invalid IDs
	invalidTestCases := []string{
		// Wrong length
		"short",
		"wayyyyyyyyyyyyyyyyyyyyyytoooooooooooooooooooooooooooooolong",

		// Invalid characters
		strings.Repeat("!", DefaultSize),
		"ID-WITH-SPECIAL-CHARS",
	}

	for _, invalidID := range invalidTestCases {
		if gen.Validate(invalidID) {
			t.Errorf("Invalid ID '%s' passed validation", invalidID)
		}
	}
}

func TestMustGenerate(t *testing.T) {
	// Test MustGenerate doesn't panic for multiple calls
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustGenerate panicked unexpectedly: %v", r)
		}
	}()

	gen := New()
	for i := 0; i < 10; i++ {
		id := gen.MustGenerate()
		if len(id) != DefaultSize {
			t.Errorf("MustGenerate produced ID of incorrect length")
		}
	}
}

func TestGlobalGenerateFunctions(t *testing.T) {
	// Test global Generate function
	id1 := Generate()
	id2 := Generate()

	if len(id1) != DefaultSize {
		t.Errorf("Global Generate() produced ID of incorrect length")
	}

	if id1 == id2 {
		t.Errorf("Global Generate() produced identical IDs")
	}

	// Test GenerateWithSize
	customSizeID := GenerateWithSize(10)
	if len(customSizeID) != 10 {
		t.Errorf("GenerateWithSize(10) did not produce 10-character ID")
	}
}

// Benchmark generate performance
func BenchmarkGenerate(b *testing.B) {
	gen := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := gen.Generate()
		if err != nil {
			b.Fatalf("Unexpected error during benchmark: %v", err)
		}
	}
}
