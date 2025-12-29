package openmeteo

import "fmt"

// ErrorType classifies the category of error that occurred during SDK operations.
type ErrorType int

const (
	// ErrorTypeValidation indicates an error due to invalid input parameters
	// (e.g., invalid coordinates, concurrent request limit exceeded).
	ErrorTypeValidation ErrorType = iota

	// ErrorTypeNetwork indicates a network or transport-level error
	// (e.g., timeout, connection failure, DNS resolution error).
	ErrorTypeNetwork

	// ErrorTypeAPI indicates an error from the Open Meteo API
	// (e.g., HTTP 4xx/5xx status codes, malformed JSON response).
	ErrorTypeAPI
)

// Error represents an error that occurred during SDK operations.
// It implements the error interface and supports error wrapping (Go 1.13+).
// Use errors.As() to extract the typed error and check the Type field programmatically.
type Error struct {
	// Type classifies the error category (Validation, Network, or API)
	Type ErrorType

	// Message is a human-readable description of the error
	Message string

	// Cause is the underlying error that caused this error (may be nil)
	Cause error
}

// Error returns a formatted error message implementing the error interface.
// If a Cause is present, it is included in the message.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause error, enabling error chain inspection with errors.Is() and errors.As().
func (e *Error) Unwrap() error {
	return e.Cause
}
