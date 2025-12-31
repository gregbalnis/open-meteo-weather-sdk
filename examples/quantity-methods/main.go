package main

import (
	"context"
	"fmt"
	"log"

	weather "github.com/gregbalnis/open-meteo-weather-sdk"
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

	// Display weather information using `QuantityOf...` methods
	fmt.Printf("Weather for %.2f°N, %.2f°E\n", w.Latitude, w.Longitude)
	fmt.Printf("Time: %s\n\n", w.Time.Format("2006-01-02 15:04 MST"))

	fmt.Println("Temperature & Humidity:")
	fmt.Printf("  Temperature:         %s\n", w.QuantityOfTemperature())
	fmt.Printf("  Feels like:          %s\n", w.QuantityOfApparentTemperature())
	fmt.Printf("  Humidity:            %s\n", w.QuantityOfRelativeHumidity())
	fmt.Println()

	fmt.Println("Wind:")
	fmt.Printf("  Speed:               %s\n", w.QuantityOfWindSpeed())
	fmt.Printf("  Direction:           %s\n", w.QuantityOfWindDirection())
	fmt.Printf("  Gusts:               %s\n", w.QuantityOfWindGusts())
	fmt.Println()

	fmt.Println("Precipitation:")
	fmt.Printf("  Total:               %s\n", w.QuantityOfPrecipitation())
	fmt.Printf("  Rain:                %s\n", w.QuantityOfRain())
	fmt.Printf("  Showers:             %s\n", w.QuantityOfShowers())
	fmt.Printf("  Snowfall:            %s\n", w.QuantityOfSnowfall())
	fmt.Println()

	fmt.Println("Atmospheric:")
	fmt.Printf("  Cloud Cover:         %s\n", w.QuantityOfCloudCover())
	fmt.Printf("  Pressure MSL:        %s\n", w.QuantityOfPressureMSL())
	fmt.Printf("  Surface Pressure:    %s\n", w.QuantityOfSurfacePressure())
}
