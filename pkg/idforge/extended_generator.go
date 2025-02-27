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

	// Check context immediately
	select {
	case <-ctx.Done():
		return "", ErrGenerationTimeout
	default:
	}

	// Prepare context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, g.config.MaxGenerationTime)
	defer cancel()

	// Collect entropy
	var entropyParts []string
	for _, provider := range g.config.Entropy {
		// Check context after each entropy collection
		select {
		case <-timeoutCtx.Done():
			return "", ErrGenerationTimeout
		default:
			entropyStr, err := provider.Provide(timeoutCtx)
			if err != nil {
				return "", err
			}
			entropyParts = append(entropyParts, entropyStr)
		}
	}

	// Determine maximum attempts more dynamically
	alphabetLen := len(g.config.Alphabet)
	maxAttempts := int(math.Min(
		math.Pow(float64(alphabetLen), float64(g.config.Size))*g.config.UniquenessPressure,
		1000, // Prevent excessive iterations
	))

	// Seed random generation with entropy
	combinedEntropy := strings.Join(entropyParts, "")
	seedBytes := []byte(combinedEntropy)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Frequently check for context cancellation
		select {
		case <-timeoutCtx.Done():
			return "", ErrGenerationTimeout
		default:
			// Generate candidate ID
			id := make([]byte, g.config.Size)
			for i := 0; i < g.config.Size; i++ {
				// Check context between each character generation
				select {
				case <-timeoutCtx.Done():
					return "", ErrGenerationTimeout
				default:
					// Use crypto/rand for secure randomness
					num, err := rand.Int(rand.Reader, big.NewInt(int64(len(g.config.Alphabet))))
					if err != nil {
						return "", err
					}

					// Incorporate entropy-based randomness
					if len(seedBytes) > 0 {
						num = new(big.Int).Add(
							num,
							big.NewInt(int64(seedBytes[i%len(seedBytes)])),
						)
						num = new(big.Int).Mod(num, big.NewInt(int64(len(g.config.Alphabet))))
					}

					id[i] = g.config.Alphabet[num.Int64()]
				}
			}

			candidateID := string(id)

			// Check for uniqueness
			if !g.generated[candidateID] {
				g.generated[candidateID] = true
				return candidateID, nil
			}
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
	possibleCombinations := math.Pow(float64(alphabetSize), float64(g.config.Size))

	// Probability of at least one collision
	probabilityOfCollision := 1 - math.Exp(
		-float64(numIDs*(numIDs-1))/(2*possibleCombinations),
	)

	// Return probability of no collisions
	return 1 - probabilityOfCollision
}
