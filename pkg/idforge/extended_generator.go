package idforge

import (
	"context"
	"crypto/rand"
	"errors"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/mrityunjay-vashisth/go-idforge/internal/entropy"
)

var (
	ErrInvalidAlphabet   = errors.New("alphabet must contain at least 2 unique characters")
	ErrInvalidSize       = errors.New("size must be positive")
	ErrGenerationTimeout = errors.New("ID generation timed out")
)

// GeneratorConfig provides advanced configuration options
type GeneratorConfig struct {
	Alphabet           string
	Size               int
	Entropy            []entropy.EntropyProvider
	MaxGenerationTime  time.Duration
	UniquenessPressure float64
}

// ExtendedGenerator provides more advanced ID generation capabilities
type ExtendedGenerator struct {
	mu        sync.Mutex
	config    GeneratorConfig
	generated map[string]bool
}

// NewExtendedGenerator creates a new generator with comprehensive configuration
func NewExtendedGenerator(opts ...func(*GeneratorConfig)) *ExtendedGenerator {
	// Default configuration
	config := GeneratorConfig{
		Alphabet:           DefaultAlphabet,
		Size:               DefaultSize,
		Entropy:            entropy.DefaultEntropyProviders(),
		MaxGenerationTime:  5 * time.Second,
		UniquenessPressure: 0.99, // 99% uniqueness guarantee
	}

	// Apply custom options
	for _, opt := range opts {
		opt(&config)
	}

	return &ExtendedGenerator{
		config:    config,
		generated: make(map[string]bool),
	}
}

// WithCustomAlphabet sets a custom character set
func WithCustomAlphabet(alphabet string) func(*GeneratorConfig) {
	return func(c *GeneratorConfig) {
		if len(alphabet) >= 2 {
			c.Alphabet = alphabet
		}
	}
}

// WithEntropyProviders allows custom entropy sources
func WithEntropyProviders(providers []entropy.EntropyProvider) func(*GeneratorConfig) {
	return func(c *GeneratorConfig) {
		if len(providers) > 0 {
			c.Entropy = providers
		}
	}
}

// Generate creates a unique identifier with advanced features
func (g *ExtendedGenerator) Generate(ctx context.Context) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Validate configuration
	if len(g.config.Alphabet) < 2 {
		return "", ErrInvalidAlphabet
	}
	if g.config.Size <= 0 {
		return "", ErrInvalidSize
	}

	// Prepare context with timeout
	ctx, cancel := context.WithTimeout(ctx, g.config.MaxGenerationTime)
	defer cancel()

	// Collect entropy
	var entropyParts []string
	for _, provider := range g.config.Entropy {
		entropyStr, err := provider.Provide(ctx)
		if err != nil {
			return "", err
		}
		entropyParts = append(entropyParts, entropyStr)
	}

	// Generate ID with uniqueness checks
	alphabetLen := big.NewInt(int64(len(g.config.Alphabet)))
	combinedEntropy := strings.Join(entropyParts, "")
	seedBytes := []byte(combinedEntropy)

	maxAttempts := int(math.Pow(float64(len(g.config.Alphabet)), float64(g.config.Size)) * g.config.UniquenessPressure)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Generate candidate ID
		id := make([]byte, g.config.Size)
		for i := 0; i < g.config.Size; i++ {
			num, err := rand.Int(rand.Reader, alphabetLen)
			if err != nil {
				return "", err
			}

			// Incorporate entropy-based randomness
			if len(seedBytes) > 0 {
				num = new(big.Int).Add(
					num,
					big.NewInt(int64(seedBytes[i%len(seedBytes)])),
				)
				num = new(big.Int).Mod(num, alphabetLen)
			}

			id[i] = g.config.Alphabet[num.Int64()]
		}

		candidateID := string(id)

		// Check for uniqueness
		if !g.generated[candidateID] {
			g.generated[candidateID] = true
			return candidateID, nil
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return "", ErrGenerationTimeout
		default:
		}
	}

	return "", ErrGenerationTimeout
}

// Validate checks if an ID meets the generator's criteria
func (g *ExtendedGenerator) Validate(id string) bool {
	if len(id) != g.config.Size {
		return false
	}

	for _, char := range id {
		if !strings.ContainsRune(g.config.Alphabet, char) {
			return false
		}
	}

	return true
}

// GetUniquenessProbability calculates the probability of generating a unique ID
func (g *ExtendedGenerator) GetUniquenessProbability(numIDs int) float64 {
	alphabetSize := len(g.config.Alphabet)
	return 1 - math.Exp(
		-float64(numIDs*(numIDs-1))/
			(2*math.Pow(float64(alphabetSize), float64(g.config.Size))),
	)
}
