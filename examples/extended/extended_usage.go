package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mrityunjay-vashisth/go-idforge/internal/entropy"
	"github.com/mrityunjay-vashisth/go-idforge/pkg/idforge"
)

func main() {
	// Create an extended generator with default configuration
	defaultExtendedGen := idforge.NewExtendedGenerator()
	defaultExtendedID, err := defaultExtendedGen.Generate(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate default extended ID: %v", err)
	}
	fmt.Println("Default Extended ID:", defaultExtendedID)

	// Create an extended generator with custom alphabet and size
	customExtendedGen := idforge.NewExtendedGenerator(
		idforge.WithCustomAlphabet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		func(cfg *idforge.GeneratorConfig) {
			cfg.Size = 16
		},
	)
	customExtendedID, err := customExtendedGen.Generate(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate custom extended ID: %v", err)
	}
	fmt.Println("Custom Extended ID:", customExtendedID)

	// Create an extended generator with custom entropy providers
	entropyExtendedGen := idforge.NewExtendedGenerator(
		idforge.WithEntropyProviders([]entropy.EntropyProvider{
			&entropy.TimestampEntropy{},
			&entropy.UUIDEntropy{},
			&entropy.RandomBytesEntropy{},
		}),
	)
	entropyExtendedID, err := entropyExtendedGen.Generate(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate entropy-based extended ID: %v", err)
	}
	fmt.Println("Entropy-based Extended ID:", entropyExtendedID)

	// Create an extended generator with a maximum generation time
	timeoutExtendedGen := idforge.NewExtendedGenerator(
		func(cfg *idforge.GeneratorConfig) {
			cfg.MaxGenerationTime = 500 * time.Millisecond
		},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	timeoutExtendedID, err := timeoutExtendedGen.Generate(ctx)
	if err != nil {
		log.Fatalf("Failed to generate timeout-based extended ID: %v", err)
	}
	fmt.Println("Timeout-based Extended ID:", timeoutExtendedID)

	// Create an extended generator with a high uniqueness pressure
	uniquenessExtendedGen := idforge.NewExtendedGenerator(
		func(cfg *idforge.GeneratorConfig) {
			cfg.UniquenessPressure = 0.9999
		},
	)
	uniquenessExtendedID, err := uniquenessExtendedGen.Generate(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate high-uniqueness extended ID: %v", err)
	}
	fmt.Println("High-uniqueness Extended ID:", uniquenessExtendedID)
}
