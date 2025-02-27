package idforge

import "github.com/mrityunjay-vashisth/go-idforge/internal/entropy"

// Option defines a function type for configuring the generator
type Option func(*Generator)

// WithAlphabet allows customizing the character set for ID generation
func WithAlphabet(alphabet string) Option {
	return func(g *Generator) {
		if len(alphabet) >= 2 {
			g.alphabet = alphabet
		}
	}
}

// WithCustomAlphabet sets a custom character set for ID generation
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

// WithSize sets the length of generated IDs
func WithSize(size int) Option {
	return func(g *Generator) {
		if size > 0 {
			g.size = size
		}
	}
}

// WithMinSize ensures a minimum ID length
func WithMinSize(minSize int) Option {
	return func(g *Generator) {
		if minSize > 0 && minSize < g.size {
			g.size = minSize
		}
	}
}

// WithMaxSize caps the maximum ID length
func WithMaxSize(maxSize int) Option {
	return func(g *Generator) {
		if maxSize > 0 && maxSize > g.size {
			g.size = maxSize
		}
	}
}
