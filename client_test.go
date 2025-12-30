package openmeteo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestGetCurrentWeather_Success tests successful weather fetching
func TestGetCurrentWeather_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		if r.URL.Query().Get("latitude") != "52.52" {
			t.Errorf("Expected latitude 52.52, got %s", r.URL.Query().Get("latitude"))
		}
		if r.URL.Query().Get("longitude") != "13.41" {
			t.Errorf("Expected longitude 13.41, got %s", r.URL.Query().Get("longitude"))
		}
		if r.URL.Query().Get("current") != "temperature_2m,relative_humidity_2m,apparent_temperature,is_day,precipitation,rain,showers,snowfall,weather_code,cloud_cover,pressure_msl,surface_pressure,wind_speed_10m,wind_direction_10m,wind_gusts_10m" {
			t.Error("Expected current=temperature_2m,relative_humidity_2m,apparent_temperature,is_day,precipitation,rain,showers,snowfall,weather_code,cloud_cover,pressure_msl,surface_pressure,wind_speed_10m,wind_direction_10m,wind_gusts_10m")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
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
		}`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	weather, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if weather == nil {
		t.Fatal("Expected weather data, got nil")
	}

	// Verify all fields
	if weather.Latitude != 52.52 {
		t.Errorf("Expected latitude 52.52, got %.2f", weather.Latitude)
	}
	if weather.Longitude != 13.41 {
		t.Errorf("Expected longitude 13.41, got %.2f", weather.Longitude)
	}
	expectedTime, _ := time.Parse("2006-01-02T15:04", "2025-12-29T10:00")
	if !weather.Time.Equal(expectedTime) {
		t.Errorf("Expected time 2025-12-29T10:00, got %s", weather.Time.Format("2006-01-02T15:04"))
	}
	if weather.Temperature != 15.3 {
		t.Errorf("Expected temperature 15.3, got %.1f", weather.Temperature)
	}
	if weather.WindSpeed != 12.5 {
		t.Errorf("Expected wind speed 12.5, got %.1f", weather.WindSpeed)
	}
	if weather.WindDirection != 270.0 {
		t.Errorf("Expected wind direction 270, got %.1f", weather.WindDirection)
	}
	if weather.WeatherCode != 3 {
		t.Errorf("Expected weather code 3, got %d", weather.WeatherCode)
	}
	if !weather.IsDay {
		t.Error("Expected IsDay to be true")
	}
	if weather.RelativeHumidity != 65.0 {
		t.Errorf("Expected relative humidity 65.0, got %.1f", weather.RelativeHumidity)
	}
	if weather.ApparentTemperature != 14.1 {
		t.Errorf("Expected apparent temperature 14.1, got %.1f", weather.ApparentTemperature)
	}
	if weather.Precipitation != 0.5 {
		t.Errorf("Expected precipitation 0.5, got %.1f", weather.Precipitation)
	}
	if weather.Rain != 0.3 {
		t.Errorf("Expected rain 0.3, got %.1f", weather.Rain)
	}
	if weather.Showers != 0.2 {
		t.Errorf("Expected showers 0.2, got %.1f", weather.Showers)
	}
	if weather.Snowfall != 0.0 {
		t.Errorf("Expected snowfall 0.0, got %.1f", weather.Snowfall)
	}
	if weather.CloudCover != 75.0 {
		t.Errorf("Expected cloud cover 75.0, got %.1f", weather.CloudCover)
	}
	if weather.PressureMSL != 1013.25 {
		t.Errorf("Expected pressure MSL 1013.25, got %.2f", weather.PressureMSL)
	}
	if weather.SurfacePressure != 1010.0 {
		t.Errorf("Expected surface pressure 1010.0, got %.1f", weather.SurfacePressure)
	}
	if weather.WindGusts != 18.0 {
		t.Errorf("Expected wind gusts 18.0, got %.1f", weather.WindGusts)
	}
}

// TestGetCurrentWeather_BoundaryCoordinates tests valid boundary coordinates
func TestGetCurrentWeather_BoundaryCoordinates(t *testing.T) {
	testCases := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"North Pole", 90.0, 0.0},
		{"South Pole", -90.0, 0.0},
		{"Date Line East", 0.0, 180.0},
		{"Date Line West", 0.0, -180.0},
		{"Northwest Corner", 90.0, -180.0},
		{"Southeast Corner", -90.0, 180.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintln(w, `{
					"latitude": 0.0,
					"longitude": 0.0,
					"current_weather": {
						"time": "2025-12-29T10:00",
						"temperature": 0.0,
						"weathercode": 0,
						"is_day": 1
					}
				}`)
			}))
			defer server.Close()

			client := NewClient(WithBaseURL(server.URL))
			_, err := client.GetCurrentWeather(context.Background(), tc.lat, tc.lon)

			if err != nil {
				t.Errorf("Expected no error for valid boundary coordinates, got %v", err)
			}
		})
	}
}

// TestGetCurrentWeather_InvalidLatitude tests invalid latitude validation
func TestGetCurrentWeather_InvalidLatitude(t *testing.T) {
	client := NewClient()

	testCases := []float64{-91.0, -100.0, 91.0, 100.0, 999.0}

	for _, lat := range testCases {
		t.Run(fmt.Sprintf("Latitude_%.0f", lat), func(t *testing.T) {
			_, err := client.GetCurrentWeather(context.Background(), lat, 0.0)

			if err == nil {
				t.Error("Expected error for invalid latitude")
			}

			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Errorf("Expected *Error, got %T", err)
			} else if apiErr.Type != ErrorTypeValidation {
				t.Errorf("Expected ErrorTypeValidation, got %v", apiErr.Type)
			}
		})
	}
}

// TestGetCurrentWeather_InvalidLongitude tests invalid longitude validation
func TestGetCurrentWeather_InvalidLongitude(t *testing.T) {
	client := NewClient()

	testCases := []float64{-181.0, -200.0, 181.0, 200.0, -999.0}

	for _, lon := range testCases {
		t.Run(fmt.Sprintf("Longitude_%.0f", lon), func(t *testing.T) {
			_, err := client.GetCurrentWeather(context.Background(), 0.0, lon)

			if err == nil {
				t.Error("Expected error for invalid longitude")
			}

			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Errorf("Expected *Error, got %T", err)
			} else if apiErr.Type != ErrorTypeValidation {
				t.Errorf("Expected ErrorTypeValidation, got %v", apiErr.Type)
			}
		})
	}
}

// TestGetCurrentWeather_ConcurrentRequests tests concurrent request handling (up to 10)
func TestGetCurrentWeather_ConcurrentRequests(t *testing.T) {
	requestCount := 0
	mu := sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()

		// Simulate some processing time
		time.Sleep(50 * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"latitude": 52.52,
			"longitude": 13.41,
			"current_weather": {
				"time": "2025-12-29T10:00",
				"temperature": 15.3,
				"weathercode": 0,
				"is_day": 1
			}
		}`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	// Launch 10 concurrent requests (should all succeed)
	var wg sync.WaitGroup
	errors := make([]error, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)
			errors[idx] = err
		}(i)
	}
	wg.Wait()

	// Verify all requests succeeded
	for i, err := range errors {
		if err != nil {
			t.Errorf("Request %d failed: %v", i, err)
		}
	}

	// Verify all requests were processed
	mu.Lock()
	if requestCount != 10 {
		t.Errorf("Expected 10 requests, got %d", requestCount)
	}
	mu.Unlock()
}

// TestGetCurrentWeather_ConcurrencyLimitExceeded tests that 11th concurrent request fails
func TestGetCurrentWeather_ConcurrencyLimitExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay to keep requests active
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"latitude": 0.0,
			"longitude": 0.0,
			"current_weather": {
				"time": "2025-12-29T10:00",
				"temperature": 0.0,
				"weathercode": 0,
				"is_day": 1
			}
		}`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	var wg sync.WaitGroup
	results := make(chan error, 15)

	// Launch 15 concurrent requests (only 10 should succeed, 5 should fail)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.GetCurrentWeather(context.Background(), 0.0, 0.0)
			results <- err
		}()
	}

	wg.Wait()
	close(results)

	successCount := 0
	failCount := 0
	for err := range results {
		if err == nil {
			successCount++
		} else {
			var apiErr *Error
			if errors.As(err, &apiErr) && apiErr.Type == ErrorTypeValidation {
				failCount++
			}
		}
	}

	if successCount != 10 {
		t.Errorf("Expected 10 successful requests, got %d", successCount)
	}
	if failCount != 5 {
		t.Errorf("Expected 5 failed requests, got %d", failCount)
	}
}

// TestGetCurrentWeather_ContextCancellation tests context cancellation
func TestGetCurrentWeather_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"latitude": 0.0, "longitude": 0.0, "current_weather": {}}`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.GetCurrentWeather(ctx, 52.52, 13.41)

	if err == nil {
		t.Error("Expected error for cancelled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

// TestGetCurrentWeather_HTTPError tests HTTP error responses
func TestGetCurrentWeather_HTTPError(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"Bad Request", http.StatusBadRequest},
		{"Internal Server Error", http.StatusInternalServerError},
		{"Service Unavailable", http.StatusServiceUnavailable},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = fmt.Fprintf(w, "Error: %s", tc.name)
			}))
			defer server.Close()

			client := NewClient(WithBaseURL(server.URL))
			_, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)

			if err == nil {
				t.Error("Expected error for non-200 status code")
			}

			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Errorf("Expected *Error, got %T", err)
			} else if apiErr.Type != ErrorTypeAPI {
				t.Errorf("Expected ErrorTypeAPI, got %v", apiErr.Type)
			}
		})
	}
}

// TestGetCurrentWeather_MalformedJSON tests malformed JSON response
func TestGetCurrentWeather_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{invalid json`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	_, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Errorf("Expected *Error, got %T", err)
	} else if apiErr.Type != ErrorTypeAPI {
		t.Errorf("Expected ErrorTypeAPI, got %v", apiErr.Type)
	}
}

// TestNewClient_DefaultConfiguration tests default client configuration
func TestNewClient_DefaultConfiguration(t *testing.T) {
	client := NewClient()

	if client.httpClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}
	if client.httpClient.Timeout != defaultTimeout {
		t.Errorf("Expected timeout %v, got %v", defaultTimeout, client.httpClient.Timeout)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", defaultBaseURL, client.baseURL)
	}
	if cap(client.semaphore) != maxConcurrent {
		t.Errorf("Expected semaphore capacity %d, got %d", maxConcurrent, cap(client.semaphore))
	}
}
