package openmeteo

import (
	"net/http"
	"testing"
	"time"
)

// TestWithTimeout tests WithTimeout option
func TestWithTimeout(t *testing.T) {
	customTimeout := 15 * time.Second
	client := NewClient(WithTimeout(customTimeout))

	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}
}

// TestWithHTTPClient tests WithHTTPClient option
func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns: 50,
		},
	}

	client := NewClient(WithHTTPClient(customClient))

	if client.httpClient != customClient {
		t.Error("Expected custom HTTP client to be used")
	}
	if client.httpClient.Timeout != 20*time.Second {
		t.Errorf("Expected timeout 20s, got %v", client.httpClient.Timeout)
	}
}

// TestWithBaseURL tests WithBaseURL option
func TestWithBaseURL(t *testing.T) {
	customURL := "https://custom-api.example.com/v2"
	client := NewClient(WithBaseURL(customURL))

	if client.baseURL != customURL {
		t.Errorf("Expected base URL %s, got %s", customURL, client.baseURL)
	}
}

// TestMultipleOptions tests combining multiple options
func TestMultipleOptions(t *testing.T) {
	customTimeout := 15 * time.Second
	customURL := "https://custom.example.com"
	customClient := &http.Client{
		Timeout: customTimeout,
	}

	client := NewClient(
		WithTimeout(customTimeout),
		WithBaseURL(customURL),
		WithHTTPClient(customClient),
	)

	if client.httpClient != customClient {
		t.Error("Expected custom HTTP client")
	}
	if client.baseURL != customURL {
		t.Errorf("Expected base URL %s, got %s", customURL, client.baseURL)
	}
	// Note: WithHTTPClient overrides timeout set by WithTimeout
	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}
}

// TestOptionOrdering tests that options are applied in order
func TestOptionOrdering(t *testing.T) {
	// First set timeout to 5s, then override with custom client having 10s
	client := NewClient(
		WithTimeout(5*time.Second),
		WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	)

	// The second option should win
	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s (last option wins), got %v", client.httpClient.Timeout)
	}
}

// TestWithTimeout_ZeroValue tests WithTimeout with zero value (no timeout)
func TestWithTimeout_ZeroValue(t *testing.T) {
	client := NewClient(WithTimeout(0))

	if client.httpClient.Timeout != 0 {
		t.Errorf("Expected timeout 0 (no timeout), got %v", client.httpClient.Timeout)
	}
}

// TestWithBaseURL_EmptyString tests WithBaseURL with empty string
func TestWithBaseURL_EmptyString(t *testing.T) {
	client := NewClient(WithBaseURL(""))

	if client.baseURL != "" {
		t.Errorf("Expected empty base URL, got %s", client.baseURL)
	}
}

// TestDefaultsWithoutOptions tests that defaults are applied when no options provided
func TestDefaultsWithoutOptions(t *testing.T) {
	client := NewClient()

	if client.httpClient.Timeout != defaultTimeout {
		t.Errorf("Expected default timeout %v, got %v", defaultTimeout, client.httpClient.Timeout)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("Expected default base URL %s, got %s", defaultBaseURL, client.baseURL)
	}
	if cap(client.semaphore) != maxConcurrent {
		t.Errorf("Expected semaphore capacity %d, got %d", maxConcurrent, cap(client.semaphore))
	}
}
