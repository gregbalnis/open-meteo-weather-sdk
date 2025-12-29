package openmeteo

import (
	"net/http"
	"time"
)

// Option is a functional option for configuring a Client.
// Options are applied during client creation via NewClient().
type Option func(*Client)

// WithTimeout sets a custom timeout for HTTP requests.
// The default timeout is 10 seconds.
//
// Example:
//
//	client := openmeteo.NewClient(openmeteo.WithTimeout(15 * time.Second))
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client for making requests.
// This is useful for customizing transport settings (e.g., proxy, TLS configuration).
//
// Example:
//
//	httpClient := &http.Client{
//	    Transport: &http.Transport{
//	        Proxy: http.ProxyFromEnvironment,
//	    },
//	}
//	client := openmeteo.NewClient(openmeteo.WithHTTPClient(httpClient))
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL for the Open Meteo API.
// This is primarily useful for testing with mock servers.
// The default base URL is https://api.open-meteo.com/v1
//
// Example:
//
//	client := openmeteo.NewClient(openmeteo.WithBaseURL("http://localhost:8080"))
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}
