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
	MaxUniqueIDs       int // New option to limit unique ID tracking
}

// ExtendedGenerator provides more advanced ID generation capabilities
type ExtendedGenerator struct {
	mu        sync.Mutex
	config    GeneratorConfig
	generated map[string]bool
	idCounter int
}

// NewExtendedGenerator creates a new generator with comprehensive configuration
func NewExtendedGenerator(opts ...func(*GeneratorConfig)) *ExtendedGenerator {
	// Default configuration
	config := GeneratorConfig{
		Alphabet:           DefaultAlphabet,
		Size:               DefaultSize,
		Entropy:            entropy.DefaultEntropyProviders(),
		MaxGenerationTime:  5 * time.Second,
		UniquenessPressure: 0.99,  // 99% uniqueness guarantee
		MaxUniqueIDs:       10000, // Limit unique ID tracking
	}

	// Apply custom options
	for _, opt := range opts {
		opt(&config)
	}

	return &ExtendedGenerator{
		config:    config,
		generated: make(map[string]bool),
		idCounter: 0,
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
	timeoutCtx, cancel := context.WithTimeout(ctx, g.config.MaxGenerationTime)
	defer cancel()

	// Efficient entropy collection with context check
	entropyParts, err := g.collectEntropy(timeoutCtx)
	if err != nil {
		return "", err
	}

	// Dynamic max attempts calculation
	alphabetLen := len(g.config.Alphabet)
	maxAttempts := calculateMaxAttempts(alphabetLen, g.config.Size, g.config.UniquenessPressure)

	// Seed random generation with entropy
	combinedEntropy := strings.Join(entropyParts, "")
	seedBytes := []byte(combinedEntropy)

	// More efficient unique ID tracking
	if g.idCounter >= g.config.MaxUniqueIDs {
		g.generated = make(map[string]bool)
		g.idCounter = 0
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Less frequent context checks
		if attempt%10 == 0 {
			select {
			case <-timeoutCtx.Done():
				return "", ErrGenerationTimeout
			default:
			}
		}

		// Generate candidate ID with optimized randomness
		candidateID := g.generateCandidateID(seedBytes)

		// Check for uniqueness
		if !g.generated[candidateID] {
			g.generated[candidateID] = true
			g.idCounter++
			return candidateID, nil
		}
	}

	return "", ErrGenerationTimeout
}

// collectEntropy efficiently gathers entropy with context management
func (g *ExtendedGenerator) collectEntropy(ctx context.Context) ([]string, error) {
	entropyParts := make([]string, 0, len(g.config.Entropy))

	for _, provider := range g.config.Entropy {
		// Occasional context check to reduce overhead
		select {
		case <-ctx.Done():
			return nil, ErrGenerationTimeout
		default:
			entropyStr, err := provider.Provide(ctx)
			if err != nil {
				return nil, err
			}
			entropyParts = append(entropyParts, entropyStr)
		}
	}

	return entropyParts, nil
}

// generateCandidateID creates an ID with enhanced randomness
func (g *ExtendedGenerator) generateCandidateID(seedBytes []byte) string {
	id := make([]byte, g.config.Size)
	alphabetLen := big.NewInt(int64(len(g.config.Alphabet)))

	for i := 0; i < g.config.Size; i++ {
		// Use crypto/rand for secure randomness
		num, _ := rand.Int(rand.Reader, alphabetLen)

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

	return string(id)
}

// Utility function to calculate max attempts dynamically
func calculateMaxAttempts(alphabetLen, size int, uniquenessPressure float64) int {
	maxAttempts := int(math.Min(
		math.Pow(float64(alphabetLen), float64(size))*uniquenessPressure,
		1000, // Prevent excessive iterations
	))
	return maxAttempts
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
