# go-idforge

[![Go Report Card](https://goreportcard.com/badge/github.com/mrityunjay-vashisth/go-idforge)](https://goreportcard.com/report/github.com/mrityunjay-vashisth/go-idforge)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

go-idforge is a robust, secure, and flexible identifier generation library for Go. It provides advanced capabilities for generating unique, cryptographically secure identifiers with extensive customization options.

## Features

- 🔐 Cryptographically secure ID generation
- 🧩 Highly customizable generation parameters
- 🔍 Advanced entropy collection
- 📊 Comprehensive validation
- 🛡️ Collision detection
- 🚀 High-performance ID generation

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
    fmt.Println(id)

    // Generate with custom configuration
    customGen := idforge.New(
        idforge.WithAlphabet("0123456789ABCDEF"),
        idforge.WithSize(8),
    )
    customID := customGen.MustGenerate()
    fmt.Println(customID)
}
```

## Extended Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/mrityunjay-vashisth/go-idforge/pkg/idforge"
)

func main() {
    // Create an extended generator with custom configuration
    extendedGen := idforge.NewExtendedGenerator(
        idforge.WithCustomAlphabet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
        func(cfg *idforge.GeneratorConfig) {
            cfg.Size = 16
            cfg.UniquenessPressure = 0.99
        },
    )

    // Generate an ID using the extended generator
    id, err := extendedGen.Generate(context.Background())
    if err != nil {
        fmt.Println("Failed to generate ID:", err)
        return
    }
    fmt.Println("Generated ID:", id)
}
```

## Error Handling

go-idforge provides detailed error messages for various scenarios. It is important to handle errors appropriately in your code. For example:

```go
id, err := generator.Generate()
if err != nil {
    // Handle the error
    log.Println("Failed to generate ID:", err)
    return
}
```

## Security Considerations

go-idforge takes security seriously and implements the following measures:

- Cryptographically secure random number generation
- Multiple entropy sources
- Timing-safe comparisons
- Comprehensive ID validation

To ensure the security of your generated IDs, follow these best practices:

- Use a sufficiently large character set and ID length
- Avoid using predictable patterns or easily guessable IDs
- Validate and sanitize IDs before using them in sensitive operations
- Keep your go-idforge version up to date to benefit from the latest security improvements

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

Thank you for using go-idforge! We hope this library serves your ID generation needs effectively and securely. If you encounter any issues or have ideas for enhancements, please let us know. Happy coding!