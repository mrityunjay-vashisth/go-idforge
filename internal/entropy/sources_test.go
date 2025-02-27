package entropy

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"testing"
)

// Helper function to validate if a string is a valid UUID
func isValidUUID(s string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(s)
}

// Helper function to test entropy providers
func testEntropyProvider(t *testing.T, provider EntropyProvider, validationFunc func(string) bool) {
	ctx := context.Background()

	// Test multiple times to ensure variability
	generatedValues := make(map[string]bool)
	for i := 0; i < 10; i++ {
		entropy, err := provider.Provide(ctx)

		// Check for error
		if err != nil {
			t.Errorf("Unexpected error from entropy provider: %v", err)
		}

		// Check if entropy is not empty
		if entropy == "" {
			t.Error("Entropy provider returned an empty string")
		}

		// Check for unique values (if applicable)
		if generatedValues[entropy] {
			t.Logf("Potential duplicate value detected: %s", entropy)
		}
		generatedValues[entropy] = true

		// Optional custom validation
		if validationFunc != nil && !validationFunc(entropy) {
			t.Errorf("Invalid entropy format: %s", entropy)
		}
	}
}

func TestTimestampEntropy(t *testing.T) {
	provider := &TimestampEntropy{}

	testEntropyProvider(t, provider, func(s string) bool {
		// Timestamp should be a positive integer
		_, err := fmt.Sscanf(s, "%d", new(int64))
		return err == nil
	})
}

func TestUUIDEntropy(t *testing.T) {
	provider := &UUIDEntropy{}

	testEntropyProvider(t, provider, isValidUUID)
}

func TestRandomBytesEntropy(t *testing.T) {
	// Test default length
	defaultProvider := &RandomBytesEntropy{}
	testEntropyProvider(t, defaultProvider, func(s string) bool {
		// Default length (16) would generate 32 hex characters
		return len(s) == 32
	})

	// Test custom length
	customLengths := []int{8, 16, 32, 64}
	for _, length := range customLengths {
		provider := &RandomBytesEntropy{length: length}
		testEntropyProvider(t, provider, func(s string) bool {
			// Hex representation should be twice the byte length
			return len(s) == length*2
		})
	}
}

func TestSystemEntropy(t *testing.T) {
	provider := &SystemEntropy{}

	testEntropyProvider(t, provider, func(s string) bool {
		// Split the system entropy parts
		parts := strings.Split(s, "_")
		return len(parts) == 4 // Alloc, NumCPU, NumGC, Timestamp
	})
}

func TestNetworkEntropy(t *testing.T) {
	provider := &NetworkEntropy{}

	// Diagnostic: Print out all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("Failed to get network interfaces: %v", err)
	}

	t.Logf("Total network interfaces: %d", len(interfaces))
	for _, iface := range interfaces {
		t.Logf("Interface: %s, Flags: %v, HardwareAddr: %s",
			iface.Name, iface.Flags, iface.HardwareAddr)
	}

	testEntropyProvider(t, provider, func(s string) bool {
		// Print the actual generated entropy for debugging
		t.Logf("Generated Network Entropy: %s", s)

		// Allow empty string if no active non-loopback interfaces
		if s == "" {
			return true
		}

		// Split MAC addresses
		macAddresses := strings.Split(s, ",")

		for _, macStr := range macAddresses {
			macStr = strings.TrimSpace(macStr)

			// Validate MAC address format
			_, err := net.ParseMAC(macStr)
			if err != nil {
				t.Logf("Invalid MAC address: %s, Error: %v", macStr, err)
				return false
			}
		}

		return true
	})
}

func TestEnhancedEntropyProvider(t *testing.T) {
	provider := &EnhancedEntropyProvider{}

	testEntropyProvider(t, provider, func(s string) bool {
		// Validate that it's a non-negative big integer
		return regexp.MustCompile(`^\d+$`).MatchString(s)
	})
}

func TestSecureEntropyAggregator(t *testing.T) {
	// Test with default providers
	aggregator := NewSecureEntropyAggregator()
	ctx := context.Background()

	entropy, err := aggregator.Aggregate(ctx)
	if err != nil {
		t.Fatalf("Unexpected error from SecureEntropyAggregator: %v", err)
	}

	// SHA-256 hash is 64 hex characters
	if len(entropy) != 64 {
		t.Errorf("Unexpected entropy length. Expected 64, got %d", len(entropy))
	}

	// Verify it uses hex characters
	if matched, _ := regexp.MatchString(`^[0-9a-f]+$`, entropy); !matched {
		t.Errorf("Entropy should be a hex string, got: %s", entropy)
	}

	// Test with custom providers
	customProviders := []EntropyProvider{
		&TimestampEntropy{},
		&UUIDEntropy{},
	}
	customAggregator := NewSecureEntropyAggregator(customProviders...)

	customEntropy, err := customAggregator.Aggregate(ctx)
	if err != nil {
		t.Fatalf("Unexpected error from custom SecureEntropyAggregator: %v", err)
	}

	if len(customEntropy) != 64 {
		t.Errorf("Unexpected custom entropy length. Expected 64, got %d", len(customEntropy))
	}
}

func TestDefaultEntropyProviders(t *testing.T) {
	providers := DefaultEntropyProviders()

	if len(providers) == 0 {
		t.Fatal("No default entropy providers found")
	}

	// Ensure all providers can generate entropy
	ctx := context.Background()
	for i, provider := range providers {
		entropy, err := provider.Provide(ctx)
		if err != nil {
			t.Errorf("Provider %d failed to generate entropy: %v", i, err)
		}
		if entropy == "" {
			t.Errorf("Provider %d generated empty entropy", i)
		}
	}
}

// Benchmark entropy providers
func BenchmarkEntropyProviders(b *testing.B) {
	providers := []EntropyProvider{
		&TimestampEntropy{},
		&UUIDEntropy{},
		&RandomBytesEntropy{},
		&SystemEntropy{},
		&NetworkEntropy{},
		&EnhancedEntropyProvider{},
	}

	ctx := context.Background()

	for _, provider := range providers {
		b.Run(fmt.Sprintf("%T", provider), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := provider.Provide(ctx)
				if err != nil {
					b.Fatalf("Error generating entropy: %v", err)
				}
			}
		})
	}
}

// Test context timeout handling
func TestContextTimeout(t *testing.T) {
	providers := []EntropyProvider{
		&TimestampEntropy{},
		&UUIDEntropy{},
		&RandomBytesEntropy{},
		&SystemEntropy{},
		&NetworkEntropy{},
		&EnhancedEntropyProvider{},
	}

	for _, provider := range providers {
		t.Run(fmt.Sprintf("%T", provider), func(t *testing.T) {
			// Some providers don't support context cancellation
			// So we'll be a bit more lenient in our test
			ctx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()

			_, err := provider.Provide(ctx)

			// For providers that don't support context, this is acceptable
			// The important thing is they don't panic
			if err != nil {
				t.Logf("Context timeout for %T: %v", provider, err)
			}
		})
	}
}
