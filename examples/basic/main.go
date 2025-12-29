package main

import (
"context"
"fmt"
"log"

weather "github.com/yourusername/open-meteo-weather-sdk"
)

func main() {
// Create a client with default configuration
client := weather.NewClient()

// Get current weather for Berlin, Germany
ctx := context.Background()
w, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
if err != nil {
log.Fatalf("Failed to fetch weather: %v", err)
}

// Display weather information
fmt.Printf("Weather for %.2f°N, %.2f°E\n", w.Latitude, w.Longitude)
fmt.Printf("Time: %s\n", w.Time.Format("2006-01-02 15:04 MST"))
fmt.Println()
fmt.Printf("Temperature:         %.1f°C\n", w.Temperature)
fmt.Printf("Feels like:          %.1f°C\n", w.ApparentTemperature)
fmt.Printf("Humidity:            %.0f%%\n", w.RelativeHumidity)
fmt.Printf("Weather Code:        %d\n", w.WeatherCode)
fmt.Printf("Day/Night:           %s\n", dayNightStr(w.IsDay))
fmt.Println()
fmt.Printf("Wind Speed:          %.1f m/s\n", w.WindSpeed)
fmt.Printf("Wind Direction:      %.0f°\n", w.WindDirection)
fmt.Printf("Wind Gusts:          %.1f m/s\n", w.WindGusts)
fmt.Println()
fmt.Printf("Precipitation:       %.1f mm\n", w.Precipitation)
fmt.Printf("Rain:                %.1f mm\n", w.Rain)
fmt.Printf("Showers:             %.1f mm\n", w.Showers)
fmt.Printf("Snowfall:            %.1f mm\n", w.Snowfall)
fmt.Println()
fmt.Printf("Cloud Cover:         %.0f%%\n", w.CloudCover)
fmt.Printf("Sea Level Pressure:  %.2f hPa\n", w.SeaLevelPressure)
fmt.Printf("Surface Pressure:    %.2f hPa\n", w.SurfacePressure)
}

func dayNightStr(isDay bool) string {
if isDay {
return "Day"
}
return "Night"
}
