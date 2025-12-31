package openmeteo

import (
	"encoding/json"
	"testing"
	"time"
)

// TestCurrentWeather_JSONUnmarshaling tests JSON unmarshaling with complete API response
func TestCurrentWeather_JSONUnmarshaling(t *testing.T) {
	jsonData := `{
		"latitude": 52.52,
		"longitude": 13.41,
		"current": {
			"time": "2025-12-29T10:00",
			"temperature_2m": 15.3,
			"wind_speed_10m": 12.5,
			"wind_direction_10m": 270.0,
			"weather_code": 3,
			"is_day": 1,
			"relative_humidity_2m": 65.0,
			"apparent_temperature": 14.1,
			"precipitation": 0.5,
			"rain": 0.3,
			"showers": 0.2,
			"snowfall": 0.0,
			"cloud_cover": 75.0,
			"pressure_msl": 1013.25,
			"surface_pressure": 1010.0,
			"wind_gusts_10m": 18.0
		}
	}`

	var resp weatherResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify all fields are correctly unmarshaled
	if resp.Latitude != 52.52 {
		t.Errorf("Expected latitude 52.52, got %.2f", resp.Latitude)
	}
	if resp.Longitude != 13.41 {
		t.Errorf("Expected longitude 13.41, got %.2f", resp.Longitude)
	}

	cw := resp.CurrentWeather
	if cw.Temperature == nil || *cw.Temperature != 15.3 {
		t.Errorf("Expected temperature 15.3, got %v", cw.Temperature)
	}
	if cw.Windspeed == nil || *cw.Windspeed != 12.5 {
		t.Errorf("Expected windspeed 12.5, got %v", cw.Windspeed)
	}
	if cw.Winddirection == nil || *cw.Winddirection != 270.0 {
		t.Errorf("Expected wind direction 270, got %v", cw.Winddirection)
	}
	if cw.Weathercode == nil || *cw.Weathercode != 3 {
		t.Errorf("Expected weather code 3, got %v", cw.Weathercode)
	}
	if cw.IsDay == nil || *cw.IsDay != 1 {
		t.Errorf("Expected is_day 1, got %v", cw.IsDay)
	}
	if cw.RelativeHumidity == nil || *cw.RelativeHumidity != 65.0 {
		t.Errorf("Expected humidity 65.0, got %v", cw.RelativeHumidity)
	}
	if cw.ApparentTemperature == nil || *cw.ApparentTemperature != 14.1 {
		t.Errorf("Expected apparent temperature 14.1, got %v", cw.ApparentTemperature)
	}
	if cw.Precipitation == nil || *cw.Precipitation != 0.5 {
		t.Errorf("Expected precipitation 0.5, got %v", cw.Precipitation)
	}
	if cw.Rain == nil || *cw.Rain != 0.3 {
		t.Errorf("Expected rain 0.3, got %v", cw.Rain)
	}
	if cw.Showers == nil || *cw.Showers != 0.2 {
		t.Errorf("Expected showers 0.2, got %v", cw.Showers)
	}
	if cw.Snowfall == nil || *cw.Snowfall != 0.0 {
		t.Errorf("Expected snowfall 0.0, got %v", cw.Snowfall)
	}
	if cw.CloudCover == nil || *cw.CloudCover != 75.0 {
		t.Errorf("Expected cloud cover 75.0, got %v", cw.CloudCover)
	}
	if cw.PressureMSL == nil || *cw.PressureMSL != 1013.25 {
		t.Errorf("Expected sea level pressure 1013.25, got %v", cw.PressureMSL)
	}
	if cw.SurfacePressure == nil || *cw.SurfacePressure != 1010.0 {
		t.Errorf("Expected surface pressure 1010.0, got %v", cw.SurfacePressure)
	}
	if cw.WindGusts == nil || *cw.WindGusts != 18.0 {
		t.Errorf("Expected wind gusts 18.0, got %v", cw.WindGusts)
	}
}

// TestCurrentWeather_JSONUnmarshalingWithNulls tests JSON unmarshaling with null/missing fields
func TestCurrentWeather_JSONUnmarshalingWithNulls(t *testing.T) {
	jsonData := `{
		"latitude": 52.52,
		"longitude": 13.41,
		"current": {
			"time": "2025-12-29T10:00",
			"temperature_2m": 15.3,
			"wind_speed_10m": null,
			"wind_direction_10m": null,
			"weather_code": 0,
			"is_day": 1
		}
	}`

	var resp weatherResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	cw := resp.CurrentWeather
	if cw.Temperature == nil || *cw.Temperature != 15.3 {
		t.Errorf("Expected temperature 15.3, got %v", cw.Temperature)
	}
	if cw.Windspeed != nil {
		t.Errorf("Expected nil windspeed, got %v", *cw.Windspeed)
	}
	if cw.Winddirection != nil {
		t.Errorf("Expected nil wind direction, got %v", *cw.Winddirection)
	}
	if cw.RelativeHumidity != nil {
		t.Errorf("Expected nil humidity, got %v", *cw.RelativeHumidity)
	}
}

// TestConvertToCurrentWeather tests conversion from API response to CurrentWeather
func TestConvertToCurrentWeather(t *testing.T) {
	c := NewClient()

	timeStr := "2025-12-29T10:00"
	temp := 15.3
	humidity := 65.0
	windspeed := 12.5
	winddirection := 270.0
	weathercode := 3
	isDay := 1

	apiResp := weatherResponse{
		Latitude:  52.52,
		Longitude: 13.41,
		CurrentWeather: currentWeatherResponse{
			Time:             &timeStr,
			Temperature:      &temp,
			RelativeHumidity: &humidity,
			Windspeed:        &windspeed,
			Winddirection:    &winddirection,
			Weathercode:      &weathercode,
			IsDay:            &isDay,
		},
	}

	weather := c.convertToCurrentWeather(apiResp)

	if weather.Latitude != 52.52 {
		t.Errorf("Expected latitude 52.52, got %.2f", weather.Latitude)
	}
	if weather.Longitude != 13.41 {
		t.Errorf("Expected longitude 13.41, got %.2f", weather.Longitude)
	}
	if weather.Temperature != 15.3 {
		t.Errorf("Expected temperature 15.3, got %.1f", weather.Temperature)
	}
	if weather.RelativeHumidity != 65.0 {
		t.Errorf("Expected humidity 65.0, got %.1f", weather.RelativeHumidity)
	}
	if weather.WindSpeed != 12.5 {
		t.Errorf("Expected wind speed 12.5, got %.1f", weather.WindSpeed)
	}
	if weather.WindDirection != 270.0 {
		t.Errorf("Expected wind direction 270.0, got %.1f", weather.WindDirection)
	}
	if weather.WeatherCode != 3 {
		t.Errorf("Expected weather code 3, got %d", weather.WeatherCode)
	}
	if !weather.IsDay {
		t.Error("Expected IsDay to be true")
	}

	expectedTime, _ := time.Parse("2006-01-02T15:04", "2025-12-29T10:00")
	if !weather.Time.Equal(expectedTime.UTC()) {
		t.Errorf("Expected time %v, got %v", expectedTime.UTC(), weather.Time)
	}
}

// TestConvertToCurrentWeather_WithNulls tests conversion with null values (should use zero values)
func TestConvertToCurrentWeather_WithNulls(t *testing.T) {
	c := NewClient()

	timeStr := "2025-12-29T10:00"
	temp := 15.3

	apiResp := weatherResponse{
		Latitude:  52.52,
		Longitude: 13.41,
		CurrentWeather: currentWeatherResponse{
			Time:        &timeStr,
			Temperature: &temp,
			// All other fields are nil
		},
	}

	weather := c.convertToCurrentWeather(apiResp)

	if weather.Temperature != 15.3 {
		t.Errorf("Expected temperature 15.3, got %.1f", weather.Temperature)
	}

	// All nil fields should be zero values
	if weather.RelativeHumidity != 0.0 {
		t.Errorf("Expected humidity 0.0, got %.1f", weather.RelativeHumidity)
	}
	if weather.WindSpeed != 0.0 {
		t.Errorf("Expected wind speed 0.0, got %.1f", weather.WindSpeed)
	}
	if weather.Precipitation != 0.0 {
		t.Errorf("Expected precipitation 0.0, got %.1f", weather.Precipitation)
	}
	if weather.WeatherCode != 0 {
		t.Errorf("Expected weather code 0, got %d", weather.WeatherCode)
	}
	if weather.IsDay {
		t.Error("Expected IsDay to be false (zero value)")
	}
}

// TestCurrentWeather_QuantityMethods tests all QuantityOf... methods
func TestCurrentWeather_QuantityMethods(t *testing.T) {
	weather := &CurrentWeather{
		Temperature:         15.3,
		ApparentTemperature: 14.1,
		RelativeHumidity:    65.0,
		Precipitation:       0.5,
		Rain:                0.3,
		Showers:             0.2,
		Snowfall:            0.0,
		CloudCover:          75.0,
		PressureMSL:         1013.25,
		SurfacePressure:     1010.0,
		WindSpeed:           12.5,
		WindDirection:       270.0,
		WindGusts:           18.0,
	}

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "QuantityOfTemperature",
			method:   weather.QuantityOfTemperature,
			expected: "15.3°C",
		},
		{
			name:     "QuantityOfApparentTemperature",
			method:   weather.QuantityOfApparentTemperature,
			expected: "14.1°C",
		},
		{
			name:     "QuantityOfRelativeHumidity",
			method:   weather.QuantityOfRelativeHumidity,
			expected: "65%",
		},
		{
			name:     "QuantityOfPrecipitation",
			method:   weather.QuantityOfPrecipitation,
			expected: "0.5 mm",
		},
		{
			name:     "QuantityOfRain",
			method:   weather.QuantityOfRain,
			expected: "0.3 mm",
		},
		{
			name:     "QuantityOfShowers",
			method:   weather.QuantityOfShowers,
			expected: "0.2 mm",
		},
		{
			name:     "QuantityOfSnowfall",
			method:   weather.QuantityOfSnowfall,
			expected: "0.0 cm",
		},
		{
			name:     "QuantityOfCloudCover",
			method:   weather.QuantityOfCloudCover,
			expected: "75%",
		},
		{
			name:     "QuantityOfPressureMSL",
			method:   weather.QuantityOfPressureMSL,
			expected: "1013.2 hPa",
		},
		{
			name:     "QuantityOfSurfacePressure",
			method:   weather.QuantityOfSurfacePressure,
			expected: "1010.0 hPa",
		},
		{
			name:     "QuantityOfWindSpeed",
			method:   weather.QuantityOfWindSpeed,
			expected: "12.5 km/h",
		},
		{
			name:     "QuantityOfWindDirection",
			method:   weather.QuantityOfWindDirection,
			expected: "270°",
		},
		{
			name:     "QuantityOfWindGusts",
			method:   weather.QuantityOfWindGusts,
			expected: "18.0 km/h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}

// TestCurrentWeather_QuantityMethods_ZeroValues tests QuantityOf... methods with zero values
func TestCurrentWeather_QuantityMethods_ZeroValues(t *testing.T) {
	weather := &CurrentWeather{}

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "QuantityOfTemperature_zero",
			method:   weather.QuantityOfTemperature,
			expected: "0.0°C",
		},
		{
			name:     "QuantityOfRelativeHumidity_zero",
			method:   weather.QuantityOfRelativeHumidity,
			expected: "0%",
		},
		{
			name:     "QuantityOfWindSpeed_zero",
			method:   weather.QuantityOfWindSpeed,
			expected: "0.0 km/h",
		},
		{
			name:     "QuantityOfPrecipitation_zero",
			method:   weather.QuantityOfPrecipitation,
			expected: "0.0 mm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}

// TestCurrentWeather_QuantityMethods_NegativeValues tests QuantityOf... methods with negative values
func TestCurrentWeather_QuantityMethods_NegativeValues(t *testing.T) {
	weather := &CurrentWeather{
		Temperature:         -10.5,
		ApparentTemperature: -15.2,
		WindDirection:       0.0,
	}

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "QuantityOfTemperature_negative",
			method:   weather.QuantityOfTemperature,
			expected: "-10.5°C",
		},
		{
			name:     "QuantityOfApparentTemperature_negative",
			method:   weather.QuantityOfApparentTemperature,
			expected: "-15.2°C",
		},
		{
			name:     "QuantityOfWindDirection_zero",
			method:   weather.QuantityOfWindDirection,
			expected: "0°",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}
