package idforge

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"sync"

	"github.com/mrityunjay-vashisth/go-idforge/internal/entropy"
)

const (
	DefaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultSize     = 21
)

type Generator struct {
	mu       sync.Mutex
	alphabet string
	size     int
	entropy  []entropy.EntropyProvider
}

func New(opts ...Option) *Generator {
	g := &Generator{
		alphabet: DefaultAlphabet,
		size:     DefaultSize,
		entropy:  entropy.DefaultEntropyProviders(),
	}

	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Generate creates a unique, secure identifier
func (g *Generator) Generate() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Collect entropy from providers
	var entropyParts []string
	ctx := context.Background()
	for _, provider := range g.entropy {
		entropyStr, err := provider.Provide(ctx)
		if err != nil {
			return "", err
		}
		entropyParts = append(entropyParts, entropyStr)
	}

	// Generate the ID using collected entropy
	id := make([]byte, g.size)
	alphabetLen := big.NewInt(int64(len(g.alphabet)))

	// Use entropy as additional randomness source
	combinedEntropy := strings.Join(entropyParts, "")
	seedBytes := []byte(combinedEntropy)

	for i := 0; i < g.size; i++ {
		// Use cryptographically secure random number generation
		num, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}

		// Add some entropy-based randomness
		if len(seedBytes) > 0 {
			num = new(big.Int).Add(
				num,
				big.NewInt(int64(seedBytes[i%len(seedBytes)])),
			)
			num = new(big.Int).Mod(num, alphabetLen)
		}

		id[i] = g.alphabet[num.Int64()]
	}

	return string(id), nil
}

// MustGenerate generates an ID, panicking on error
func (g *Generator) MustGenerate() string {
	id, err := g.Generate()
	if err != nil {
		panic(err)
	}
	return id
}

// Validate checks if an ID meets the generator's criteria
func (g *Generator) Validate(id string) bool {
	if len(id) != g.size {
		return false
	}

	for _, char := range id {
		if !strings.ContainsRune(g.alphabet, char) {
			return false
		}
	}

	return true
}

// Quick generation functions for convenience
func Generate() string {
	return New().MustGenerate()
}

func GenerateWithSize(size int) string {
	return New(WithSize(size)).MustGenerate()
}
