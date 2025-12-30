package openmeteo

import (
	"errors"
	"fmt"
	"testing"
)

// TestError_ErrorMethod tests Error.Error() formatting
func TestError_ErrorMethod(t *testing.T) {
	testCases := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "Without cause",
			err: &Error{
				Type:    ErrorTypeValidation,
				Message: "invalid latitude",
			},
			expected: "invalid latitude",
		},
		{
			name: "With cause",
			err: &Error{
				Type:    ErrorTypeNetwork,
				Message: "network timeout",
				Cause:   fmt.Errorf("connection refused"),
			},
			expected: "network timeout: connection refused",
		},
		{
			name: "API error with cause",
			err: &Error{
				Type:    ErrorTypeAPI,
				Message: "failed to parse JSON",
				Cause:   fmt.Errorf("unexpected end of JSON input"),
			},
			expected: "failed to parse JSON: unexpected end of JSON input",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.err.Error()
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

// TestError_Unwrap tests Error.Unwrap() for error chain inspection
func TestError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := &Error{
		Type:    ErrorTypeNetwork,
		Message: "network error",
		Cause:   cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}

// TestError_UnwrapNil tests Error.Unwrap() when Cause is nil
func TestError_UnwrapNil(t *testing.T) {
	err := &Error{
		Type:    ErrorTypeValidation,
		Message: "validation error",
		Cause:   nil,
	}

	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Errorf("Expected nil, got %v", unwrapped)
	}
}

// TestError_ErrorsIs tests compatibility with errors.Is()
func TestError_ErrorsIs(t *testing.T) {
	rootCause := fmt.Errorf("root cause")
	wrappedErr := fmt.Errorf("wrapped: %w", rootCause)
	apiErr := &Error{
		Type:    ErrorTypeAPI,
		Message: "API error",
		Cause:   wrappedErr,
	}

	// Test errors.Is with the root cause
	if !errors.Is(apiErr, rootCause) {
		t.Error("errors.Is should find root cause in chain")
	}

	// Test errors.Is with wrapped error
	if !errors.Is(apiErr, wrappedErr) {
		t.Error("errors.Is should find wrapped error in chain")
	}

	// Test errors.Is with unrelated error
	unrelatedErr := fmt.Errorf("unrelated")
	if errors.Is(apiErr, unrelatedErr) {
		t.Error("errors.Is should not match unrelated error")
	}
}

// TestError_ErrorsAs tests compatibility with errors.As()
func TestError_ErrorsAs(t *testing.T) {
	originalErr := &Error{
		Type:    ErrorTypeValidation,
		Message: "validation failed",
	}

	// Wrap the error
	wrappedErr := fmt.Errorf("operation failed: %w", originalErr)

	// Test errors.As extraction
	var targetErr *Error
	if !errors.As(wrappedErr, &targetErr) {
		t.Fatal("errors.As should extract *Error from wrapped error")
	}

	if targetErr.Type != ErrorTypeValidation {
		t.Errorf("Expected ErrorTypeValidation, got %v", targetErr.Type)
	}
	if targetErr.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got %q", targetErr.Message)
	}
}

// TestErrorType_Values tests ErrorType constant values
func TestErrorType_Values(t *testing.T) {
	// Ensure error types have distinct values
	types := map[ErrorType]string{
		ErrorTypeValidation: "Validation",
		ErrorTypeNetwork:    "Network",
		ErrorTypeAPI:        "API",
	}

	seen := make(map[ErrorType]bool)
	for typ, name := range types {
		if seen[typ] {
			t.Errorf("Duplicate ErrorType value for %s", name)
		}
		seen[typ] = true
	}

	if len(seen) != 3 {
		t.Errorf("Expected 3 distinct ErrorType values, got %d", len(seen))
	}
}

// TestError_TypeSwitch tests programmatic error type checking
func TestError_TypeSwitch(t *testing.T) {
	testCases := []struct {
		name         string
		err          *Error
		expectedType ErrorType
	}{
		{
			name: "Validation error",
			err: &Error{
				Type:    ErrorTypeValidation,
				Message: "invalid input",
			},
			expectedType: ErrorTypeValidation,
		},
		{
			name: "Network error",
			err: &Error{
				Type:    ErrorTypeNetwork,
				Message: "timeout",
			},
			expectedType: ErrorTypeNetwork,
		},
		{
			name: "API error",
			err: &Error{
				Type:    ErrorTypeAPI,
				Message: "server error",
			},
			expectedType: ErrorTypeAPI,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var apiErr *Error
			if !errors.As(tc.err, &apiErr) {
				t.Fatal("errors.As failed")
			}

			switch apiErr.Type {
			case ErrorTypeValidation:
				if tc.expectedType != ErrorTypeValidation {
					t.Error("Unexpected validation error")
				}
			case ErrorTypeNetwork:
				if tc.expectedType != ErrorTypeNetwork {
					t.Error("Unexpected network error")
				}
			case ErrorTypeAPI:
				if tc.expectedType != ErrorTypeAPI {
					t.Error("Unexpected API error")
				}
			default:
				t.Errorf("Unknown error type: %v", apiErr.Type)
			}
		})
	}
}
