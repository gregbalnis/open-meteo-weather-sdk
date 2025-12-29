# Quick Start Guide: Open Meteo Current Weather SDK

**Date**: 2025-12-28

## Installation

```bash
go get github.com/yourusername/open-meteo-weather-sdk
```

**Requirements**: Go 1.21 or later

---

## Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/yourusername/open-meteo-weather-sdk"
)

func main() {
    // Create a client
    client := openmeteo.NewClient()
    
    // Get current weather for Berlin (52.52°N, 13.41°E)
    weather, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)
    if err != nil {
        log.Fatal(err)
    }
    
    // Display results
    fmt.Printf("Location: %.2f°N, %.2f°E\n", weather.Latitude, weather.Longitude)
    fmt.Printf("Temperature: %.1f°C\n", weather.Temperature)
    fmt.Printf("Feels like: %.1f°C\n", weather.ApparentTemperature)
    fmt.Printf("Humidity: %.0f%%\n", weather.RelativeHumidity)
    fmt.Printf("Wind Speed: %.1f m/s\n", weather.WindSpeed)
    fmt.Printf("Wind Direction: %.0f°\n", weather.WindDirection)
}
```

**Output**:
```
Location: 52.52°N, 13.41°E
Temperature: 15.3°C
Feels like: 14.1°C
Humidity: 65%
Wind Speed: 12.5 m/s
Wind Direction: 270°
```

---

## Custom Configuration

### Custom Timeout

```go
import "time"

// Set 15-second timeout instead of default 10 seconds
client := openmeteo.NewClient(
    openmeteo.WithTimeout(15 * time.Second),
)
```

### Custom HTTP Client

```go
import "net/http"

// Use custom HTTP client with proxy
httpClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyFromEnvironment,
    },
}

client := openmeteo.NewClient(
    openmeteo.WithHTTPClient(httpClient),
)
```

### Multiple Options

```go
client := openmeteo.NewClient(
    openmeteo.WithTimeout(15 * time.Second),
    openmeteo.WithHTTPClient(customHTTPClient),
)
```

---

## Error Handling

### Basic Error Checking

```go
weather, err := client.GetCurrentWeather(ctx, lat, lon)
if err != nil {
    log.Printf("Failed to get weather: %v", err)
    return
}
```

### Typed Error Handling

```go
import "errors"

weather, err := client.GetCurrentWeather(ctx, lat, lon)
if err != nil {
    var apiErr *openmeteo.Error
    if errors.As(err, &apiErr) {
        switch apiErr.Type {
        case openmeteo.ErrorTypeValidation:
            // Invalid coordinates or rate limit
            log.Printf("Validation error: %v", err)
        case openmeteo.ErrorTypeNetwork:
            // Network timeout or connection failure
            log.Printf("Network error (retrying...): %v", err)
        case openmeteo.ErrorTypeAPI:
            // API returned error or malformed response
            log.Printf("API error: %v", err)
        }
    }
    return
}
```

---

## Context and Cancellation

### With Timeout

```go
// Request will be cancelled after 5 seconds
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

weather, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Request timed out")
    }
    return
}
```

### Manual Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())

// Cancel request after user input
go func() {
    fmt.Scanln() // Wait for Enter key
    cancel()
}()

weather, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
if err != nil {
    if errors.Is(err, context.Canceled) {
        log.Println("Request cancelled by user")
    }
    return
}
```

---

## Concurrent Requests

### Multiple Locations

```go
import "sync"

func main() {
    client := openmeteo.NewClient()
    
    locations := []struct{
        name string
        lat, lon float64
    }{
        {"Berlin", 52.52, 13.41},
        {"New York", 40.71, -74.01},
        {"Tokyo", 35.68, 139.65},
        {"Sydney", -33.87, 151.21},
    }
    
    var wg sync.WaitGroup
    ctx := context.Background()
    
    for _, loc := range locations {
        wg.Add(1)
        go func(name string, lat, lon float64) {
            defer wg.Done()
            
            weather, err := client.GetCurrentWeather(ctx, lat, lon)
            if err != nil {
                log.Printf("%s: error: %v", name, err)
                return
            }
            
            fmt.Printf("%s: %.1f°C, Wind: %.1f m/s\n", 
                name, weather.Temperature, weather.WindSpeed)
        }(loc.name, loc.lat, loc.lon)
    }
    
    wg.Wait()
}
```

**Note**: The SDK enforces a maximum of 10 concurrent requests per client instance.

---

## Working with Weather Data

### All Available Fields

```go
weather, _ := client.GetCurrentWeather(ctx, lat, lon)

// Location
fmt.Printf("Latitude: %.2f°\n", weather.Latitude)
fmt.Printf("Longitude: %.2f°\n", weather.Longitude)
fmt.Printf("Time: %s\n", weather.Time.Format("2006-01-02 15:04:05 MST"))

// Temperature
fmt.Printf("Temperature: %.1f°C\n", weather.Temperature)
fmt.Printf("Apparent Temperature: %.1f°C\n", weather.ApparentTemperature)

// Humidity
fmt.Printf("Relative Humidity: %.0f%%\n", weather.RelativeHumidity)

// Daylight
if weather.IsDay {
    fmt.Println("Condition: Daytime")
} else {
    fmt.Println("Condition: Nighttime")
}

// Precipitation
fmt.Printf("Precipitation: %.1fmm\n", weather.Precipitation)
fmt.Printf("Rain: %.1fmm\n", weather.Rain)
fmt.Printf("Showers: %.1fmm\n", weather.Showers)
fmt.Printf("Snowfall: %.1fmm\n", weather.Snowfall)

// Sky
fmt.Printf("Weather Code: %d\n", weather.WeatherCode)
fmt.Printf("Cloud Cover: %.0f%%\n", weather.CloudCover)

// Pressure
fmt.Printf("Sea Level Pressure: %.2f hPa\n", weather.SeaLevelPressure)
fmt.Printf("Surface Pressure: %.2f hPa\n", weather.SurfacePressure)

// Wind
fmt.Printf("Wind Speed: %.1f m/s\n", weather.WindSpeed)
fmt.Printf("Wind Direction: %.0f°\n", weather.WindDirection)
fmt.Printf("Wind Gusts: %.1f m/s\n", weather.WindGusts)
```

### Weather Code Interpretation

```go
func describeWeather(code int) string {
    switch {
    case code == 0:
        return "Clear sky"
    case code <= 3:
        return "Partly cloudy"
    case code <= 48:
        return "Foggy"
    case code <= 57:
        return "Drizzle"
    case code <= 67:
        return "Rain"
    case code <= 77:
        return "Snow"
    case code <= 82:
        return "Rain showers"
    case code <= 86:
        return "Snow showers"
    default:
        return "Thunderstorm"
    }
}

weather, _ := client.GetCurrentWeather(ctx, lat, lon)
fmt.Printf("Conditions: %s\n", describeWeather(weather.WeatherCode))
```

---

## Testing Your Code

### Mock Server for Tests

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestMyWeatherFunc(t *testing.T) {
    // Create mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{
            "latitude": 52.52,
            "longitude": 13.41,
            "current": {
                "time": "2025-12-28T10:00",
                "temperature_2m": 15.3,
                "relativehumidity_2m": 65.0,
                "windspeed_10m": 12.5
            }
        }`))
    }))
    defer server.Close()
    
    // Use mock server URL
    client := openmeteo.NewClient(
        openmeteo.WithBaseURL(server.URL),
    )
    
    // Test your code
    weather, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if weather.Temperature != 15.3 {
        t.Errorf("expected temperature 15.3, got %.1f", weather.Temperature)
    }
}
```

---

## Best Practices

### 1. Reuse Client Instances

✅ **Good**: Reuse client for connection pooling
```go
client := openmeteo.NewClient()
// Use client for multiple requests
weather1, _ := client.GetCurrentWeather(ctx, lat1, lon1)
weather2, _ := client.GetCurrentWeather(ctx, lat2, lon2)
```

❌ **Bad**: Creating new client for each request
```go
// Inefficient - creates new HTTP client each time
weather1, _ := openmeteo.NewClient().GetCurrentWeather(ctx, lat1, lon1)
weather2, _ := openmeteo.NewClient().GetCurrentWeather(ctx, lat2, lon2)
```

### 2. Always Use Context

✅ **Good**: Pass context for cancellation/timeout
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
weather, err := client.GetCurrentWeather(ctx, lat, lon)
```

❌ **Bad**: Using empty context
```go
// No timeout protection
weather, err := client.GetCurrentWeather(context.Background(), lat, lon)
```

### 3. Handle Errors Properly

✅ **Good**: Check error types and handle appropriately
```go
weather, err := client.GetCurrentWeather(ctx, lat, lon)
if err != nil {
    var apiErr *openmeteo.Error
    if errors.As(err, &apiErr) {
        if apiErr.Type == openmeteo.ErrorTypeNetwork {
            // Retry logic
        }
    }
    return err
}
```

❌ **Bad**: Ignoring errors
```go
weather, _ := client.GetCurrentWeather(ctx, lat, lon)
// Proceeding without checking error
```

### 4. Validate Coordinates

✅ **Good**: Validate before calling API
```go
if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
    return fmt.Errorf("invalid coordinates")
}
weather, err := client.GetCurrentWeather(ctx, lat, lon)
```

The SDK validates coordinates, but pre-validation improves error messages in your application.

---

## Common Patterns

### Weather Dashboard

```go
func refreshWeather(client *openmeteo.Client, lat, lon float64) {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    
    for {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        weather, err := client.GetCurrentWeather(ctx, lat, lon)
        cancel()
        
        if err != nil {
            log.Printf("Error fetching weather: %v", err)
            continue
        }
        
        displayWeather(weather)
        <-ticker.C
    }
}
```

### Weather Comparison

```go
func compareLocations(client *openmeteo.Client, locations map[string][2]float64) {
    results := make(map[string]*openmeteo.CurrentWeather)
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    ctx := context.Background()
    
    for name, coords := range locations {
        wg.Add(1)
        go func(n string, lat, lon float64) {
            defer wg.Done()
            weather, err := client.GetCurrentWeather(ctx, lat, lon)
            if err != nil {
                log.Printf("%s: %v", n, err)
                return
            }
            mu.Lock()
            results[n] = weather
            mu.Unlock()
        }(name, coords[0], coords[1])
    }
    
    wg.Wait()
    
    // Find warmest location
    var warmest string
    var maxTemp float64 = -999
    for name, weather := range results {
        if weather.Temperature > maxTemp {
            maxTemp = weather.Temperature
            warmest = name
        }
    }
    
    fmt.Printf("Warmest: %s (%.1f°C)\n", warmest, maxTemp)
}
```

---

## Troubleshooting

### Issue: "concurrent request limit exceeded"

**Cause**: More than 10 simultaneous requests from one client

**Solution**: 
- Use goroutine limiting (e.g., worker pool)
- Create multiple client instances
- Queue requests

### Issue: "network timeout"

**Cause**: Request exceeded timeout (default 10s)

**Solution**:
- Increase timeout: `WithTimeout(15 * time.Second)`
- Check network connectivity
- Check Open Meteo API status

### Issue: "invalid coordinates"

**Cause**: Coordinates outside valid range

**Solution**: Ensure latitude is -90 to 90, longitude is -180 to 180

---

## Next Steps

- Read the [API documentation](contracts/api.md) for complete API reference
- Check [examples/](../examples/) for more code samples
- Review the [data model](data-model.md) for field details
- See [open-meteo-api.md](contracts/open-meteo-api.md) for API contract details

---

## Support

- **Issues**: https://github.com/yourusername/open-meteo-weather-sdk/issues
- **Open Meteo API**: https://open-meteo.com/en/docs
- **Go Documentation**: https://pkg.go.dev/github.com/yourusername/open-meteo-weather-sdk
