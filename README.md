# Open Meteo Weather SDK

A Go SDK for fetching current weather data from the [Open Meteo API](https://open-meteo.com/).

## Features

- ✅ Fetch current weather data by coordinates (latitude/longitude)
- ✅ Thread-safe client with concurrency control (max 10 simultaneous requests)
- ✅ Typed error handling (validation, network, API errors)
- ✅ Configurable timeouts and HTTP client
- ✅ Zero external dependencies (stdlib only)
- ✅ 80%+ test coverage

## Installation

```bash
go get github.com/yourusername/open-meteo-weather-sdk
```

**Requirements**: Go 1.21 or later

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    weather "github.com/yourusername/open-meteo-weather-sdk"
)

func main() {
    // Create a client
    client := weather.NewClient()
    
    // Get current weather for Berlin
    w, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Temperature: %.1f°C\n", w.Temperature)
    fmt.Printf("Wind Speed: %.1f m/s\n", w.WindSpeed)
}
```

## Usage

### Custom Configuration

```go
import "time"

// Custom timeout (default: 10s)
client := weather.NewClient(
    weather.WithTimeout(15 * time.Second),
)

// Custom HTTP client
httpClient := &http.Client{
    Transport: myTransport,
}
client := weather.NewClient(
    weather.WithHTTPClient(httpClient),
)
```

### Error Handling

```go
w, err := client.GetCurrentWeather(ctx, lat, lon)
if err != nil {
    var apiErr *weather.Error
    if errors.As(err, &apiErr) {
        switch apiErr.Type {
        case weather.ErrorTypeValidation:
            // Invalid input (coordinates, rate limit)
        case weather.ErrorTypeNetwork:
            // Network failure
        case weather.ErrorTypeAPI:
            // API error
        }
    }
}
```

## API Reference

See [GoDoc](https://pkg.go.dev/github.com/gregbalnis/open-meteo-weather-sdk) for complete API documentation.

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Check coverage (requires 80%)
make coverage

# Clean artifacts
make clean
```

## License

This project is licensed under the terms of the MIT open source license. Please refer to the [LICENSE](https://github.com/gregbalnis/open-meteo-weather-sdk/blob/main/LICENSE) file for the full terms.

## Contributing

Contributions welcome! Please open an issue or pull request.
