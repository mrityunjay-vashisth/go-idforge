# go-idforge

[![Go Report Card](https://goreportcard.com/badge/github.com/mrityunjay-vashisth/go-idforge)](https://goreportcard.com/report/github.com/mrityunjay-vashisth/go-idforge)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

go-idforge is a robust, secure, and flexible identifier generation library for Go. It provides advanced capabilities for generating unique, cryptographically secure identifiers with extensive customization options.

## Features

- 🔐 Cryptographically secure ID generation
- 🧩 Highly customizable generation parameters
- 🔍 Advanced multi-source entropy collection
- 📊 Comprehensive validation and collision detection
- ⏱️ Context-aware generation with timeout support
- 🚀 High-performance, concurrent-safe ID generation

## Installation

To install go-idforge, use the following command:

```bash
go get github.com/mrityunjay-vashisth/go-idforge
```

## Basic Usage

```go
package main

import (
    "fmt"
    "github.com/mrityunjay-vashisth/go-idforge/pkg/idforge"
)

func main() {
    // Generate a default ID
    id := idforge.Generate()
    fmt.Println("Default ID:", id)

    // Generate with custom size
    shortID := idforge.GenerateWithSize(10)
    fmt.Println("Short ID:", shortID)

    // Generate with custom configuration
    customGen := idforge.New(
        idforge.WithAlphabet("0123456789ABCDEF"),
        idforge.WithSize(8),
    )
    customID := customGen.MustGenerate()
    fmt.Println("Custom ID:", customID)
}
```

## Extended Usage

The library offers an extended generator with advanced features:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/mrityunjay-vashisth/go-idforge/pkg/idforge"
    "github.com/mrityunjay-vashisth/go-idforge/internal/entropy"
)

func main() {
    // Create an extended generator with custom configuration
    extendedGen := idforge.NewExtendedGenerator(
        idforge.WithCustomAlphabet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
        func(cfg *idforge.GeneratorConfig) {
            cfg.Size = 16
            cfg.UniquenessPressure = 0.99
            cfg.MaxGenerationTime = 2 * time.Second
        },
    )

    // Generate with context support
    ctx := context.Background()
    id, err := extendedGen.Generate(ctx)
    if err != nil {
        log.Fatalf("Failed to generate ID: %v", err)
    }
    fmt.Println("Generated ID:", id)
    
    // Calculate uniqueness probability
    probability := extendedGen.GetUniquenessProbability(1000)
    fmt.Printf("Probability of no collisions across 1000 IDs: %.8f\n", probability)
}
```

## Secure Token Generation

Besides ID generation, the library provides utilities for secure token generation:

```go
// Generate a cryptographically secure token
token, err := idforge.GenerateSecureToken(32)
if err != nil {
    log.Fatalf("Failed to generate token: %v", err)
}
fmt.Println("Secure Token:", token)

// Generate a token with panic on error
panicToken := idforge.MustGenerateSecureToken(24)
fmt.Println("Must Generate Token:", panicToken)
```

## Advanced Entropy Collection

The library uses multiple entropy sources to ensure high-quality randomness:

- Timestamp-based entropy
- UUID-based entropy
- Cryptographically secure random bytes
- System information (memory usage, CPU count, GC stats)
- Network interface information
- Enhanced entropy with aggregation and hashing

You can customize entropy providers:

```go
customGen := idforge.NewExtendedGenerator(
    idforge.WithEntropyProviders([]entropy.EntropyProvider{
        &entropy.TimestampEntropy{},
        &entropy.UUIDEntropy{},
        &entropy.RandomBytesEntropy{},
        &entropy.SystemEntropy{},
    }),
)
```

## Customization Options

### Basic Generator Options

- `WithAlphabet(string)`: Define custom character set for IDs
- `WithSize(int)`: Set exact ID length
- `WithMinSize(int)`: Ensure minimum ID length
- `WithMaxSize(int)`: Cap maximum ID length

### Extended Generator Options

- `WithCustomAlphabet(string)`: Define custom character set
- `WithEntropyProviders([]entropy.EntropyProvider)`: Custom entropy sources
- Custom configuration via function:
  ```go
  func(cfg *idforge.GeneratorConfig) {
      cfg.Size = 16
      cfg.MaxGenerationTime = 2 * time.Second
      cfg.UniquenessPressure = 0.99
      cfg.MaxUniqueIDs = 10000
  }
  ```

## ID Validation

Both generators provide methods to validate IDs:

```go
// Basic validation
isValid := idforge.IsValidID(id, "0123456789ABCDEF", 8)

// Generator-specific validation
gen := idforge.New(idforge.WithAlphabet("0123456789ABCDEF"))
isValid := gen.Validate(id)
```

## Error Handling

The library provides comprehensive error handling:

```go
id, err := generator.Generate()
if err != nil {
    switch {
    case errors.Is(err, idforge.ErrInvalidAlphabet):
        // Handle invalid alphabet
    case errors.Is(err, idforge.ErrInvalidSize):
        // Handle invalid size
    case errors.Is(err, idforge.ErrGenerationTimeout):
        // Handle timeout
    default:
        // Handle other errors
    }
}
```

## Performance Considerations

For high-volume ID generation:

1. Reuse generator instances rather than creating new ones
2. Use appropriate alphabet and size for your use case
3. Set realistic `MaxGenerationTime` for your application
4. Adjust `UniquenessPressure` based on uniqueness requirements
5. Set appropriate `MaxUniqueIDs` to limit memory consumption

## Security Considerations

go-idforge takes security seriously and implements the following measures:

- Cryptographically secure random number generation with `crypto/rand`
- Multiple entropy sources for maximum unpredictability
- Secure token generation with Base32 encoding
- Comprehensive ID validation

To ensure the security of your generated IDs:

- Use a sufficiently large character set and ID length
- Avoid using predictable patterns or easily guessable IDs
- Validate and sanitize IDs before using them in sensitive operations
- Keep your go-idforge version up to date

## Contributing

Contributions to go-idforge are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the [GitHub repository](https://github.com/mrityunjay-vashisth/go-idforge).

When contributing, please follow the existing code style and conventions. Make sure to write unit tests for any new features or changes.

## License

go-idforge is released under the [MIT License](https://opensource.org/licenses/MIT).

## Contact

If you have any questions, suggestions, or feedback, please feel free to reach out:

- Email: [mrityunjayvashisth@gmail.com]
- GitHub: [https://github.com/mrityunjay-vashisth]

---

Thank you for using go-idforge! We hope this library serves your ID generation needs effectively and securely.