# SDK Public API Contract

**Date**: 2025-12-28
**Feature**: Open Meteo Current Weather SDK
**Language**: Go

## Overview

This document defines the public API contract for the Open Meteo Current Weather SDK. All signatures follow Go conventions and are designed for simplicity, type safety, and ease of use.

## Package Declaration

```go
package openmeteo
```

**Import Path**: `github.com/gregbalnis/open-meteo-weather-sdk`

## Public API

### Client Type

```go
// Client is the main entry point for interacting with the Open Meteo API.
// It is safe for concurrent use by multiple goroutines.
type Client struct {
    // contains filtered or unexported fields
}
```

### Constructor

```go
// NewClient creates a new Open Meteo API client with optional configuration.
// The client is ready to use immediately and is safe for concurrent access.
//
// Default configuration:
//   - HTTP timeout: 10 seconds
//   - Base URL: https://api.open-meteo.com/v1
//   - Max concurrent requests: 10
//
// Example:
//   client := openmeteo.NewClient()
//   client := openmeteo.NewClient(openmeteo.WithTimeout(15 * time.Second))
func NewClient(opts ...Option) *Client
```

**Parameters**:
- `opts ...Option` - Zero or more configuration options

**Returns**:
- `*Client` - Configured client instance, never nil

**Thread Safety**: Safe to call from multiple goroutines

**Example**:
```go
// Default client
client := openmeteo.NewClient()

// Custom timeout
client := openmeteo.NewClient(
    openmeteo.WithTimeout(15 * time.Second),
)

// Custom HTTP client
httpClient := &http.Client{/* ... */}
client := openmeteo.NewClient(
    openmeteo.WithHTTPClient(httpClient),
)
```

---

### Primary Method

```go
// GetCurrentWeather retrieves current weather conditions for the specified coordinates.
// All 15 weather parameters are returned in a single request.
//
// Coordinates must be in the range:
//   - Latitude: -90.0 to 90.0 (inclusive)
//   - Longitude: -180.0 to 180.0 (inclusive)
//
// The method enforces a maximum of 10 concurrent requests per Client instance.
// If this limit is exceeded, ErrorTypeValidation is returned immediately.
//
// Network operations respect the configured timeout (default: 10 seconds).
// Context cancellation is supported for early termination.
//
// When the API returns null for any weather parameter, zero values are used:
//   - Numeric fields: 0.0
//   - Boolean fields: false
//
// All temperature values are in Celsius, wind speeds in m/s, precipitation in mm,
// and pressure in hPa.
//
// Example:
//   ctx := context.Background()
//   weather, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
//   if err != nil {
//       // Handle error
//   }
//   fmt.Printf("Temperature: %.1f°C\n", weather.Temperature)
func (c *Client) GetCurrentWeather(ctx context.Context, latitude, longitude float64) (*CurrentWeather, error)
```

**Parameters**:
- `ctx context.Context` - Context for cancellation and deadlines
- `latitude float64` - Latitude in decimal degrees (-90 to 90)
- `longitude float64` - Longitude in decimal degrees (-180 to 180)

**Returns**:
- `*CurrentWeather` - Weather data for the location (nil on error)
- `error` - Error if request fails (nil on success)

**Error Types**:
- `ErrorTypeValidation` - Invalid coordinates or concurrent limit exceeded
- `ErrorTypeNetwork` - Network timeout, connection failure, or transport error
- `ErrorTypeAPI` - API returned error status or malformed response

**Thread Safety**: Safe to call from multiple goroutines concurrently (up to 10 simultaneous requests)

**Performance**: Typical response time 100-500ms depending on network latency

**Example**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

weather, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
if err != nil {
    var apiErr *openmeteo.Error
    if errors.As(err, &apiErr) {
        switch apiErr.Type {
        case openmeteo.ErrorTypeValidation:
            log.Printf("Invalid input: %v", err)
        case openmeteo.ErrorTypeNetwork:
            log.Printf("Network error: %v", err)
        case openmeteo.ErrorTypeAPI:
            log.Printf("API error: %v", err)
        }
    }
    return
}

fmt.Printf("Temperature: %.1f°C\n", weather.Temperature)
fmt.Printf("Wind Speed: %.1f m/s\n", weather.WindSpeed)
```

---

### Configuration Options

```go
// Option configures a Client instance.
type Option func(*Client)
```

#### WithTimeout

```go
// WithTimeout sets a custom timeout for HTTP requests.
// The timeout applies to the entire request-response cycle.
//
// Default: 10 seconds
//
// Example:
//   client := openmeteo.NewClient(openmeteo.WithTimeout(15 * time.Second))
func WithTimeout(timeout time.Duration) Option
```

**Parameters**:
- `timeout time.Duration` - Request timeout (must be > 0)

**Panics**: If timeout <= 0

---

#### WithHTTPClient

```go
// WithHTTPClient sets a custom HTTP client for requests.
// This is useful for configuring proxies, custom transports, or TLS settings.
//
// The provided client's Timeout field is ignored; use WithTimeout instead.
//
// Example:
//   httpClient := &http.Client{
//       Transport: &http.Transport{
//           Proxy: http.ProxyFromEnvironment,
//       },
//   }
//   client := openmeteo.NewClient(openmeteo.WithHTTPClient(httpClient))
func WithHTTPClient(httpClient *http.Client) Option
```

**Parameters**:
- `httpClient *http.Client` - Custom HTTP client (must not be nil)

**Panics**: If httpClient is nil

---

#### WithBaseURL

```go
// WithBaseURL sets a custom base URL for API requests.
// This is primarily useful for testing with mock servers.
//
// The URL must use HTTPS scheme for production use.
//
// Default: https://api.open-meteo.com/v1
//
// Example:
//   client := openmeteo.NewClient(openmeteo.WithBaseURL("https://test.example.com"))
func WithBaseURL(baseURL string) Option
```

**Parameters**:
- `baseURL string` - Base URL for API requests (must be valid URL)

**Panics**: If baseURL is not a valid URL

---

### Data Types

```go
// CurrentWeather represents current weather conditions at a specific location.
// All fields use metric units:
//   - Temperature: Celsius (°C)
//   - Wind speed: meters per second (m/s)
//   - Precipitation: millimeters (mm)
//   - Pressure: hectopascals (hPa)
//   - Percentages: 0-100 (%)
type CurrentWeather struct {
    // Location
    Latitude  float64   // Latitude in decimal degrees
    Longitude float64   // Longitude in decimal degrees
    Time      time.Time // Observation time (UTC)
    
    // Temperature
    Temperature         float64 // Air temperature at 2m height (°C)
    ApparentTemperature float64 // Feels-like temperature (°C)
    
    // Humidity
    RelativeHumidity float64 // Relative humidity at 2m (%)
    
    // Daylight
    IsDay bool // True if daytime, false if night
    
    // Precipitation
    Precipitation float64 // Total precipitation (mm)
    Rain          float64 // Liquid rain amount (mm)
    Showers       float64 // Shower precipitation (mm)
    Snowfall      float64 // Snowfall water equivalent (mm)
    
    // Sky conditions
    WeatherCode int     // WMO weather code (0-99)
    CloudCover  float64 // Total cloud cover (%)
    
    // Pressure
    SeaLevelPressure float64 // Pressure at sea level (hPa)
    SurfacePressure  float64 // Pressure at surface (hPa)
    
    // Wind
    WindSpeed     float64 // Wind speed at 10m height (m/s)
    WindDirection float64 // Wind direction at 10m (degrees, 0=N)
    WindGusts     float64 // Wind gust speed at 10m (m/s)
}
```

**Notes**:
- All fields are exported (public)
- Zero values indicate no measurement (e.g., 0.0mm rain means no rain)
- Time is always in UTC timezone
- WeatherCode follows WMO standard (https://www.nodc.noaa.gov/archive/arc0021/0002199/1.1/data/0-data/HTML/WMO-CODE/WMO4677.HTM)

---

### Error Types

```go
// Error represents an error from the SDK with type classification.
// It implements the error interface and supports error unwrapping.
type Error struct {
    Type    ErrorType // Error classification
    Message string    // Human-readable error message
    Cause   error     // Underlying error (may be nil)
}
```

**Methods**:

```go
// Error returns the error message.
// Implements the error interface.
func (e *Error) Error() string

// Unwrap returns the underlying error for error chain inspection.
// Returns nil if there is no underlying error.
func (e *Error) Unwrap() error
```

---

```go
// ErrorType classifies errors into categories.
type ErrorType int

const (
    // ErrorTypeValidation indicates invalid input parameters.
    // Examples: invalid coordinates, concurrent limit exceeded
    ErrorTypeValidation ErrorType = iota
    
    // ErrorTypeNetwork indicates network or transport failures.
    // Examples: timeout, connection refused, DNS failure
    ErrorTypeNetwork
    
    // ErrorTypeAPI indicates API-specific errors.
    // Examples: HTTP 500, malformed JSON response
    ErrorTypeAPI
)
```

**Usage**:
```go
var apiErr *openmeteo.Error
if errors.As(err, &apiErr) {
    switch apiErr.Type {
    case openmeteo.ErrorTypeValidation:
        // Handle validation error
    case openmeteo.ErrorTypeNetwork:
        // Handle network error
    case openmeteo.ErrorTypeAPI:
        // Handle API error
    }
}
```

---

## Contract Guarantees

### Stability

1. **Semantic Versioning**: Breaking changes will increment major version
2. **Backward Compatibility**: Adding fields to `CurrentWeather` is non-breaking
3. **Error Types**: New error types may be added without breaking existing error handling

### Thread Safety

1. `Client` instances are safe for concurrent use
2. `CurrentWeather` instances are immutable after creation
3. `Error` instances are immutable

### Performance

1. Connection pooling is automatic (reuses HTTP connections)
2. No goroutine leaks; all resources cleaned up on context cancellation
3. Zero allocations in steady state (excluding JSON parsing)

### Error Handling

1. All errors implement standard `error` interface
2. Error chains support `errors.Is()` and `errors.As()`
3. Error messages are human-readable and actionable

### Data Integrity

1. Coordinate validation prevents invalid API requests
2. JSON parsing validates field types
3. Null values from API are converted to zero values (documented behavior)

---

## Examples

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/gregbalnis/open-meteo-weather-sdk"
)

func main() {
    client := openmeteo.NewClient()
    
    weather, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Temperature: %.1f°C\n", weather.Temperature)
    fmt.Printf("Wind Speed: %.1f m/s\n", weather.WindSpeed)
}
```

### With Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

weather, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
// ...
```

### Error Handling

```go
weather, err := client.GetCurrentWeather(ctx, lat, lon)
if err != nil {
    var apiErr *openmeteo.Error
    if errors.As(err, &apiErr) {
        switch apiErr.Type {
        case openmeteo.ErrorTypeValidation:
            log.Printf("Invalid input: %v", err)
        case openmeteo.ErrorTypeNetwork:
            log.Printf("Network error, retrying...: %v", err)
        case openmeteo.ErrorTypeAPI:
            log.Printf("API error: %v", err)
        }
    }
    return
}
```

### Concurrent Usage

```go
var wg sync.WaitGroup
locations := []struct{ lat, lon float64 }{
    {52.52, 13.41},  // Berlin
    {40.71, -74.01}, // New York
    {35.68, 139.65}, // Tokyo
}

for _, loc := range locations {
    wg.Add(1)
    go func(lat, lon float64) {
        defer wg.Done()
        weather, err := client.GetCurrentWeather(ctx, lat, lon)
        if err != nil {
            log.Printf("Error: %v", err)
            return
        }
        fmt.Printf("%.2f,%.2f: %.1f°C\n", lat, lon, weather.Temperature)
    }(loc.lat, loc.lon)
}

wg.Wait()
```

---

## Versioning

**Initial Version**: v0.1.0

**Stability**: Beta (pre-1.0)
- API may change before 1.0.0
- Breaking changes will be documented in CHANGELOG
- After 1.0.0: Semantic Versioning strictly followed
