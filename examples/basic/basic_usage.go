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

	// Create a custom generator
	customGen := idforge.New(
		idforge.WithAlphabet("0123456789ABCDEF"),
		idforge.WithSize(12),
	)

	// Generate ID using custom generator
	customID := customGen.MustGenerate()
	fmt.Println("Custom ID:", customID)

	// Create an extended generator with custom configuration
	extendedGen := idforge.NewExtendedGenerator(
		idforge.WithCustomAlphabet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		func(cfg *idforge.GeneratorConfig) {
			cfg.Size = 16
		},
	)

	// Generate ID using extended generator
	extendedID, err := extendedGen.Generate(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate extended ID: %v", err)
	}
	fmt.Println("Extended ID:", extendedID)

	// Generate a secure token
	secureToken, err := idforge.GenerateSecureToken(32)
	if err != nil {
		log.Fatalf("Failed to generate secure token: %v", err)
	}
	fmt.Println("Secure Token:", secureToken)

	// Generate a secure token using MustGenerateSecureToken
	mustSecureToken := idforge.MustGenerateSecureToken(24)
	fmt.Println("Must Secure Token:", mustSecureToken)
}
