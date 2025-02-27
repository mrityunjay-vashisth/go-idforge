package entropy

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EntropyProvider defines an interface for generating entropy
type EntropyProvider interface {
	Provide(ctx context.Context) (string, error)
}

// TimestampEntropy provides entropy based on current timestamp
type TimestampEntropy struct{}

func (t *TimestampEntropy) Provide(ctx context.Context) (string, error) {
	return fmt.Sprintf("%d", time.Now().UnixNano()), nil
}

// UUIDEntropy generates entropy using UUID
type UUIDEntropy struct{}

func (u *UUIDEntropy) Provide(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}

// RandomBytesEntropy generates entropy from cryptographically secure random bytes
type RandomBytesEntropy struct {
	length int
}

func (r *RandomBytesEntropy) Provide(ctx context.Context) (string, error) {
	if r.length == 0 {
		r.length = 16 // Default length
	}

	b := make([]byte, r.length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

// SystemEntropy collects entropy from system-related information
type SystemEntropy struct{}

func (s *SystemEntropy) Provide(ctx context.Context) (string, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Collect various system-related values
	return fmt.Sprintf(
		"%d_%d_%d_%d",
		memStats.Alloc,
		runtime.NumCPU(),
		memStats.NumGC,
		time.Now().UnixNano(),
	), nil
}

// NetworkEntropy generates entropy from network interfaces
type NetworkEntropy struct{}

func (n *NetworkEntropy) Provide(ctx context.Context) (string, error) {
	// Get network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// Collect MAC addresses from non-loopback, up interfaces
	var macAddresses []string
	for _, iface := range interfaces {
		// Check if interface is up and not a loopback
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			// Only add if HardwareAddr is not empty
			if len(iface.HardwareAddr) > 0 {
				macAddresses = append(macAddresses, iface.HardwareAddr.String())
			}
		}
	}

	// If no MAC addresses found, return empty string
	if len(macAddresses) == 0 {
		return "", nil
	}

	// Join MAC addresses
	return strings.Join(macAddresses, ","), nil
}

// EnhancedEntropyProvider adds more sophisticated entropy generation
type EnhancedEntropyProvider struct {
	mu        sync.Mutex
	lastValue *big.Int
}

func (e *EnhancedEntropyProvider) Provide(ctx context.Context) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Collect multiple entropy sources
	sources := [][]byte{
		// Timestamp with nanosecond precision
		binary.BigEndian.AppendUint64(nil, uint64(time.Now().UnixNano())),

		// UUID as bytes
		[]byte(uuid.New().String()),

		// Memory statistics
		func() []byte {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			return binary.BigEndian.AppendUint64(nil, memStats.Alloc)
		}(),

		// Goroutine ID (somewhat unique)
		[]byte(fmt.Sprintf("%d", runtime.NumGoroutine())),
	}

	// Combine sources using SHA-256
	hash := sha256.New()
	for _, source := range sources {
		hash.Write(source)
	}

	// If a previous value exists, incorporate it for additional randomness
	if e.lastValue != nil {
		hash.Write(e.lastValue.Bytes())
	}

	// Generate a new big integer from the hash
	hashBytes := hash.Sum(nil)
	newValue := new(big.Int).SetBytes(hashBytes)

	// Store the last generated value
	e.lastValue = newValue

	return newValue.String(), nil
}

// SecureEntropyAggregator combines multiple entropy sources with additional security
type SecureEntropyAggregator struct {
	providers []EntropyProvider
}

func NewSecureEntropyAggregator(providers ...EntropyProvider) *SecureEntropyAggregator {
	// Add default enhanced entropy if no providers specified
	if len(providers) == 0 {
		providers = []EntropyProvider{
			&EnhancedEntropyProvider{},
			&SystemEntropy{},
			&UUIDEntropy{},
		}
	}
	return &SecureEntropyAggregator{providers: providers}
}

func (s *SecureEntropyAggregator) Aggregate(ctx context.Context) (string, error) {
	var entropyParts []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errs []error

	for _, provider := range s.providers {
		wg.Add(1)
		go func(p EntropyProvider) {
			defer wg.Done()
			entropy, err := p.Provide(ctx)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			mu.Lock()
			entropyParts = append(entropyParts, entropy)
			mu.Unlock()
		}(provider)
	}

	wg.Wait()

	if len(errs) > 0 {
		return "", fmt.Errorf("entropy collection errors: %v", errs)
	}

	// Hash the combined entropy for additional security
	hash := sha256.New()
	for _, part := range entropyParts {
		hash.Write([]byte(part))
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// DefaultEntropyProviders returns a set of standard entropy sources
func DefaultEntropyProviders() []EntropyProvider {
	return []EntropyProvider{
		&TimestampEntropy{},
		&UUIDEntropy{},
		&RandomBytesEntropy{length: 16},
		&SystemEntropy{},
		&EnhancedEntropyProvider{},
	}
}
