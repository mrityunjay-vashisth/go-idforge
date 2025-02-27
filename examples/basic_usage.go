package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mrityunjay-vashisth/go-idforge/pkg/idforge"
)

func main() {
	// Basic ID Generation
	basicID := idforge.Generate()
	fmt.Println("Basic ID:", basicID)

	// Generate with custom size
	shortID := idforge.GenerateWithSize(10)
	fmt.Println("Short ID:", shortID)

	// Create extended generator with custom configuration
	extendedGen := idforge.NewExtendedGenerator(
		idforge.WithCustomAlphabet("0123456789ABCDEF"),
		func(cfg *idforge.GeneratorConfig) {
			cfg.Size = 8
		},
	)

	// Generate ID using extended generator
	customID, err := extendedGen.Generate(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate custom ID: %v", err)
	}
	fmt.Println("Custom Hex ID:", customID)

	// Create a validator
	validator := idforge.NewIDValidator(
		idforge.WithMinLength(8),
		idforge.WithMaxLength(16),
		idforge.AddForbiddenPattern(`sequential`),
	)

	// Validate the generated ID
	if err := validator.Validate(customID); err != nil {
		log.Printf("ID validation error: %v", err)
	} else {
		fmt.Println("ID passed validation")
	}

	// Generate complexity report
	complexityReport := idforge.GenerateComplexityReport(customID)
	fmt.Println("ID Complexity Report:")
	for k, v := range complexityReport {
		fmt.Printf("%s: %v\n", k, v)
	}

	// Calculate uniqueness probability
	uniquenessProbability := extendedGen.GetUniquenessProbability(1000)
	fmt.Printf("Uniqueness Probability for 1000 IDs: %.4f%%\n",
		uniquenessProbability*100)

	// Secure comparison example
	anotherID := idforge.Generate()
	areEqual := idforge.SecureCompare(customID, anotherID)
	fmt.Printf("IDs are equal: %v\n", areEqual)

	// Sanitize ID example
	sanitizedID := idforge.SanitizeID(customID+"!@#", "0123456789ABCDEF")
	fmt.Println("Sanitized ID:", sanitizedID)
}
