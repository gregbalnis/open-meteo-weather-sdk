package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://api.open-meteo.com/v1"
	defaultTimeout = 10 * time.Second
	maxConcurrent  = 10
)

// Client is the main SDK entry point for making weather data requests.
// It is thread-safe and can be shared across multiple goroutines.
// Create instances using NewClient().
type Client struct {
	// httpClient is the HTTP client used for making API requests
	httpClient *http.Client

	// baseURL is the base URL for the Open Meteo API
	baseURL string

	// semaphore controls concurrent request limits (max 10 simultaneous requests)
	semaphore chan struct{}
}

// NewClient creates a new Open Meteo API client with default configuration.
// The default configuration includes:
//   - 10-second request timeout
//   - Base URL: https://api.open-meteo.com/v1
//   - Maximum 10 concurrent requests
//
// Use functional options to customize the client behavior:
//
//	client := openmeteo.NewClient(
//	    openmeteo.WithTimeout(15 * time.Second),
//	    openmeteo.WithHTTPClient(customHTTPClient),
//	)
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL:   defaultBaseURL,
		semaphore: make(chan struct{}, maxConcurrent),
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetCurrentWeather fetches current weather data for the specified geographic coordinates.
// It returns all 15 weather parameters including temperature, humidity, wind, precipitation, etc.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - latitude: Latitude in degrees (-90 to 90)
//   - longitude: Longitude in degrees (-180 to 180)
//
// Returns:
//   - *CurrentWeather: Complete weather data snapshot
//   - error: Typed error (use errors.As to extract *Error and check Type field)
//
// Example:
//
//	weather, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
//	if err != nil {
//	    var apiErr *openmeteo.Error
//	    if errors.As(err, &apiErr) {
//	        switch apiErr.Type {
//	        case openmeteo.ErrorTypeValidation:
//	            // Handle validation error
//	        case openmeteo.ErrorTypeNetwork:
//	            // Handle network error
//	        case openmeteo.ErrorTypeAPI:
//	            // Handle API error
//	        }
//	    }
//	    return err
//	}
func (c *Client) GetCurrentWeather(ctx context.Context, latitude, longitude float64) (*CurrentWeather, error) {
	// Validate coordinates
	if latitude < -90 || latitude > 90 {
		return nil, &Error{
			Type:    ErrorTypeValidation,
			Message: fmt.Sprintf("invalid latitude: %.2f (must be between -90 and 90)", latitude),
		}
	}
	if longitude < -180 || longitude > 180 {
		return nil, &Error{
			Type:    ErrorTypeValidation,
			Message: fmt.Sprintf("invalid longitude: %.2f (must be between -180 and 180)", longitude),
		}
	}

	// Acquire semaphore (concurrency control)
	select {
	case c.semaphore <- struct{}{}:
		defer func() { <-c.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, &Error{
			Type:    ErrorTypeValidation,
			Message: fmt.Sprintf("concurrent request limit exceeded (%d)", maxConcurrent),
		}
	}

	// Build request URL
	reqURL, err := c.buildRequestURL(latitude, longitude)
	if err != nil {
		return nil, &Error{
			Type:    ErrorTypeValidation,
			Message: "failed to build request URL",
			Cause:   err,
		}
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, &Error{
			Type:    ErrorTypeNetwork,
			Message: "failed to create HTTP request",
			Cause:   err,
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &Error{
			Type:    ErrorTypeNetwork,
			Message: "failed to execute HTTP request",
			Cause:   err,
		}
	}
	defer func() { _ = resp.Body.Close() }()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &Error{
			Type:    ErrorTypeAPI,
			Message: fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)),
		}
	}

	// Parse JSON response
	var apiResp weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, &Error{
			Type:    ErrorTypeAPI,
			Message: "failed to parse JSON response",
			Cause:   err,
		}
	}

	// Convert to CurrentWeather
	weather := c.convertToCurrentWeather(apiResp)
	return weather, nil
}

// buildRequestURL constructs the API request URL with query parameters
func (c *Client) buildRequestURL(latitude, longitude float64) (string, error) {
	u, err := url.Parse(c.baseURL + "/forecast")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	q.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	q.Set("current_weather", "true")
	q.Set("temperature_unit", "celsius")
	q.Set("windspeed_unit", "ms")
	q.Set("precipitation_unit", "mm")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// convertToCurrentWeather converts the internal API response to the public CurrentWeather type.
// Null values from the API are converted to zero values.
func (c *Client) convertToCurrentWeather(apiResp weatherResponse) *CurrentWeather {
	cw := &CurrentWeather{
		Latitude:  apiResp.Latitude,
		Longitude: apiResp.Longitude,
	}

	// Parse time
	if apiResp.CurrentWeather.Time != nil {
		if t, err := time.Parse("2006-01-02T15:04", *apiResp.CurrentWeather.Time); err == nil {
			cw.Time = t.UTC()
		}
	}

	// Copy fields with null handling (use zero values for nil pointers)
	if apiResp.CurrentWeather.Temperature != nil {
		cw.Temperature = *apiResp.CurrentWeather.Temperature
	}
	if apiResp.CurrentWeather.RelativeHumidity != nil {
		cw.RelativeHumidity = *apiResp.CurrentWeather.RelativeHumidity
	}
	if apiResp.CurrentWeather.ApparentTemperature != nil {
		cw.ApparentTemperature = *apiResp.CurrentWeather.ApparentTemperature
	}
	if apiResp.CurrentWeather.IsDay != nil {
		cw.IsDay = *apiResp.CurrentWeather.IsDay == 1
	}
	if apiResp.CurrentWeather.Precipitation != nil {
		cw.Precipitation = *apiResp.CurrentWeather.Precipitation
	}
	if apiResp.CurrentWeather.Rain != nil {
		cw.Rain = *apiResp.CurrentWeather.Rain
	}
	if apiResp.CurrentWeather.Showers != nil {
		cw.Showers = *apiResp.CurrentWeather.Showers
	}
	if apiResp.CurrentWeather.Snowfall != nil {
		cw.Snowfall = *apiResp.CurrentWeather.Snowfall
	}
	if apiResp.CurrentWeather.Weathercode != nil {
		cw.WeatherCode = *apiResp.CurrentWeather.Weathercode
	}
	if apiResp.CurrentWeather.CloudCover != nil {
		cw.CloudCover = *apiResp.CurrentWeather.CloudCover
	}
	if apiResp.CurrentWeather.PressureMsl != nil {
		cw.SeaLevelPressure = *apiResp.CurrentWeather.PressureMsl
	}
	if apiResp.CurrentWeather.SurfacePressure != nil {
		cw.SurfacePressure = *apiResp.CurrentWeather.SurfacePressure
	}
	if apiResp.CurrentWeather.Windspeed != nil {
		cw.WindSpeed = *apiResp.CurrentWeather.Windspeed
	}
	if apiResp.CurrentWeather.Winddirection != nil {
		cw.WindDirection = *apiResp.CurrentWeather.Winddirection
	}
	if apiResp.CurrentWeather.WindGusts != nil {
		cw.WindGusts = *apiResp.CurrentWeather.WindGusts
	}

	return cw
}
