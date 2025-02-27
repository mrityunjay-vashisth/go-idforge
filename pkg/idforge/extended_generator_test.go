package idforge

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/mrityunjay-vashisth/go-idforge/internal/entropy"
)

func TestNewExtendedGenerator(t *testing.T) {
	// Test default configuration
	defaultGen := NewExtendedGenerator()

	if defaultGen.config.Alphabet != DefaultAlphabet {
		t.Errorf("Default alphabet incorrect. Expected %s, got %s",
			DefaultAlphabet, defaultGen.config.Alphabet)
	}

	if defaultGen.config.Size != DefaultSize {
		t.Errorf("Default size incorrect. Expected %d, got %d",
			DefaultSize, defaultGen.config.Size)
	}

	if defaultGen.config.MaxGenerationTime != 5*time.Second {
		t.Errorf("Default max generation time incorrect. Expected 5s, got %v",
			defaultGen.config.MaxGenerationTime)
	}
}

func TestExtendedGeneratorWithCustomOptions(t *testing.T) {
	// Custom alphabet
	customAlphabet := "0123456789"

	// Custom entropy providers
	customProviders := []entropy.EntropyProvider{
		&entropy.TimestampEntropy{},
	}

	gen := NewExtendedGenerator(
		WithCustomAlphabet(customAlphabet),
		WithEntropyProviders(customProviders),
		func(cfg *GeneratorConfig) {
			cfg.Size = 10
			cfg.MaxGenerationTime = 2 * time.Second
			cfg.UniquenessPressure = 0.95
		},
	)

	// Verify custom configuration
	if gen.config.Alphabet != customAlphabet {
		t.Errorf("Custom alphabet not set correctly. Got %s", gen.config.Alphabet)
	}

	if gen.config.Size != 10 {
		t.Errorf("Custom size not set correctly. Got %d", gen.config.Size)
	}

	if gen.config.MaxGenerationTime != 2*time.Second {
		t.Errorf("Custom max generation time not set correctly. Got %v",
			gen.config.MaxGenerationTime)
	}

	if gen.config.UniquenessPressure != 0.95 {
		t.Errorf("Custom uniqueness pressure not set correctly. Got %f",
			gen.config.UniquenessPressure)
	}
}

func TestExtendedGeneratorGenerate(t *testing.T) {
	// Use a more lenient configuration for testing
	gen := NewExtendedGenerator(
		func(cfg *GeneratorConfig) {
			cfg.MaxGenerationTime = 30 * time.Second
			cfg.UniquenessPressure = 0.5 // Reduce uniqueness pressure
		},
	)
	ctx := context.Background()

	// Generate multiple IDs to test various properties
	generatedIDs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := gen.Generate(ctx)

		// Check for generation error
		if err != nil {
			t.Fatalf("Unexpected error generating ID on attempt %d: %v", i, err)
		}

		// Check ID length
		if len(id) != gen.config.Size {
			t.Errorf("Generated ID length incorrect. Expected %d, got %d",
				gen.config.Size, len(id))
		}

		// Check characters are from the defined alphabet
		for _, char := range id {
			if !strings.ContainsRune(gen.config.Alphabet, char) {
				t.Errorf("ID contains character not in alphabet: %c", char)
			}
		}

		// Check for uniqueness
		if generatedIDs[id] {
			t.Logf("Potential duplicate ID generated: %s", id)
		}
		generatedIDs[id] = true
	}
}

func TestExtendedGeneratorUniquenessProbability(t *testing.T) {
	// Test with different configurations
	testCases := []struct {
		name           string
		alphabetSize   int
		idSize         int
		expectedNumIDs int
	}{
		{
			name:           "Small Alphabet, Small ID",
			alphabetSize:   10,
			idSize:         5,
			expectedNumIDs: 50,
		},
		{
			name:           "Large Alphabet, Large ID",
			alphabetSize:   62,
			idSize:         21,
			expectedNumIDs: 1000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create generator with specific alphabet
			alphabet := DefaultAlphabet[:tc.alphabetSize]
			gen := NewExtendedGenerator(
				WithCustomAlphabet(alphabet),
				func(cfg *GeneratorConfig) {
					cfg.Size = tc.idSize
				},
			)

			// Calculate uniqueness probability
			probability := gen.GetUniquenessProbability(tc.expectedNumIDs)

			// Log detailed information
			t.Logf("Alphabet size: %d, ID size: %d, Expected IDs: %d, Probability: %f",
				tc.alphabetSize, tc.idSize, tc.expectedNumIDs, probability)

			// Ensure probability is a valid float between 0 and 1
			if math.IsNaN(probability) || probability < 0 || probability > 1 {
				t.Errorf("Invalid probability calculation. Got %f", probability)
			}
		})
	}
}

// Benchmark generator performance
func BenchmarkExtendedGenerator(b *testing.B) {
	gen := NewExtendedGenerator()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(ctx)
		if err != nil {
			b.Fatalf("Unexpected error during benchmark: %v", err)
		}
	}
}
