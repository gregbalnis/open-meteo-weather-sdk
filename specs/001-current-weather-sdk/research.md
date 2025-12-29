# Research: Open Meteo Current Weather SDK

**Date**: 2025-12-28
**Feature**: Current Weather API Integration

## Overview

This document consolidates research findings for implementing a Go SDK that wraps the Open Meteo Current Weather API. All technical unknowns from the specification have been resolved with concrete decisions backed by research.

## Research Questions & Findings

### 1. Open Meteo API Investigation

**Question**: What is the exact API endpoint structure, required parameters, and response format for current weather data?

**Research**: Analyzed Open Meteo API documentation at https://open-meteo.com/en/docs

**Decision**: 
- **Endpoint**: `https://api.open-meteo.com/v1/forecast`
- **Required Parameters**: `latitude`, `longitude`, `current_weather=true`
- **Optional Parameters**: `temperature_unit`, `windspeed_unit`, `precipitation_unit` (though we're using metric only per clarifications)
- **Response Format**: JSON with structure:
  ```json
  {
    "latitude": 52.52,
    "longitude": 13.41,
    "current_weather": {
      "temperature": 15.3,
      "windspeed": 12.5,
      "winddirection": 270,
      "weathercode": 3,
      "is_day": 1,
      "time": "2025-12-28T10:00"
    }
  }
  ```

**Rationale**: The API is RESTful, free (no API key required), and returns consistent JSON. The `current_weather=true` parameter enables fetching all 15 weather parameters in a single request.

**Alternatives Considered**: 
- Using hourly or daily endpoints instead → Rejected because spec requires current conditions only
- Separate requests per parameter → Rejected due to inefficiency; single request provides all data

---

### 2. Go HTTP Client Best Practices

**Question**: What is the recommended pattern for building HTTP clients in Go that are thread-safe, testable, and performant?

**Research**: Reviewed Go standard library patterns, popular SDK implementations (AWS SDK Go, Google Cloud Go), and community best practices.

**Decision**: Use the standard `net/http` package with:
- Custom `http.Client` with configurable timeout (default 10s)
- Reusable client instance (connection pooling)
- Context support for cancellation and deadlines
- Interface-based design for testability

**Rationale**: Go's `net/http` is production-ready, well-tested, and efficient. Connection pooling is automatic when reusing client instances. Custom timeouts prevent hanging requests.

**Pattern**:
```go
type Client struct {
    httpClient *http.Client
    baseURL    string
    semaphore  chan struct{} // For concurrency limiting
}

func NewClient(opts ...Option) *Client {
    c := &Client{
        httpClient: &http.Client{Timeout: 10 * time.Second},
        baseURL:    "https://api.open-meteo.com/v1",
        semaphore:  make(chan struct{}, 10), // Max 10 concurrent
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}
```

**Alternatives Considered**:
- Third-party HTTP libraries (Resty, Gentleman) → Rejected to minimize dependencies per constitution
- Raw `http.DefaultClient` → Rejected because it lacks timeout configuration

---

### 3. Error Handling Patterns in Go

**Question**: How should we structure error types to distinguish between validation, network, and API errors while following Go idioms?

**Research**: Studied Go 1.13+ error wrapping, `errors.Is/As`, and SDK error patterns from established libraries.

**Decision**: Create custom error types implementing `error` interface:
```go
type ErrorType int

const (
    ErrorTypeValidation ErrorType = iota
    ErrorTypeNetwork
    ErrorTypeAPI
)

type Error struct {
    Type    ErrorType
    Message string
    Cause   error
}

func (e *Error) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *Error) Unwrap() error {
    return e.Cause
}
```

**Rationale**: This pattern allows callers to check error types using `errors.As()`, provides context, and follows Go 1.13+ error wrapping conventions. Human-readable messages satisfy UX consistency principle.

**Alternatives Considered**:
- String-based errors → Rejected; not type-safe or programmatically distinguishable
- Sentinel errors → Rejected; less flexible for varied error scenarios
- Error codes (int) → Rejected; less idiomatic than type-based errors in Go

---

### 4. Concurrency Control

**Question**: What is the best way to limit concurrent requests to 10 per SDK instance while maintaining thread safety?

**Research**: Investigated Go concurrency patterns: semaphores using buffered channels, `sync.WaitGroup`, worker pools.

**Decision**: Buffered channel semaphore pattern:
```go
type Client struct {
    semaphore chan struct{} // Size 10
}

func (c *Client) GetCurrentWeather(ctx context.Context, lat, lon float64) (*CurrentWeather, error) {
    select {
    case c.semaphore <- struct{}{}: // Acquire
        defer func() { <-c.semaphore }() // Release
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        return nil, &Error{Type: ErrorTypeValidation, Message: "concurrent request limit exceeded (10)"}
    }
    // ... perform request
}
```

**Rationale**: Buffered channels provide efficient, lock-free concurrency control. The `select` with `default` enables fail-fast behavior when limit is reached. Context integration allows cancellation.

**Alternatives Considered**:
- `sync.Mutex` with counter → Rejected; channel-based semaphore is more idiomatic
- Rate limiter (e.g., golang.org/x/time/rate) → Rejected; spec requires instant failure, not throttling
- No limit → Rejected per clarifications; must enforce 10 concurrent max

---

### 5. Testing Strategy

**Question**: How do we achieve 80% coverage with deterministic, fast tests that don't hit the real API?

**Research**: Reviewed Go testing best practices, httptest package, and mocking strategies.

**Decision**: Three-tier testing approach:
1. **Unit tests**: Test logic with mocked HTTP responses using `httptest.Server`
2. **Integration tests**: Optional real API tests (marked with build tag `//go:build integration`)
3. **Contract tests**: Validate request/response structure matches Open Meteo API contract

**Pattern**:
```go
func TestGetCurrentWeather(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintln(w, `{"latitude":52.52,"current_weather":{"temperature":15.3}}`)
    }))
    defer server.Close()
    
    client := NewClient(WithBaseURL(server.URL))
    weather, err := client.GetCurrentWeather(context.Background(), 52.52, 13.41)
    // assertions...
}
```

**Rationale**: `httptest` is standard library, requires no external dependencies, provides full HTTP mock capability. Build tags separate integration tests from fast unit tests.

**Alternatives Considered**:
- Interface mocking libraries (testify/mock) → Rejected; `httptest` is simpler for HTTP clients
- Recording/playback (go-vcr) → Rejected; adds dependency and complexity
- Only integration tests → Rejected; slow and unreliable (network dependency)

---

### 6. Go Module Structure

**Question**: What is the best package layout for a single-purpose SDK following Go conventions?

**Research**: Studied Go project layout guidelines, standard library patterns, and popular SDK structures.

**Decision**: Flat package structure at repository root:
```
open-meteo-weather-sdk/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── Makefile
├── .github/
│   └── workflows/
│       └── ci.yml
├── client.go          # Main SDK client
├── client_test.go
├── weather.go         # Weather data types
├── weather_test.go
├── errors.go          # Error types
├── errors_test.go
└── examples/
    └── basic/
        └── main.go
```

**Rationale**: For single-purpose libraries, a flat structure at the root is most idiomatic. Package name matches repository name. Avoids unnecessary nesting (no `pkg/` directory needed). Tests colocated with code per constitution.

**Alternatives Considered**:
- `pkg/` subdirectory → Rejected; unnecessary for single-package library
- Internal packages structure → Rejected; no private implementation needed for MVP
- Separate packages per concern → Rejected; premature abstraction for ~500 LOC project

---

### 7. CI/CD with GitHub Actions

**Question**: How do we configure GitHub Actions for Go with linting, testing, and coverage enforcement?

**Research**: Reviewed golangci-lint-action documentation, Go testing coverage tools, and GitHub Actions best practices.

**Decision**: Single workflow file `.github/workflows/ci.yml`:
```yaml
name: CI
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest
  
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $coverage%"
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "Coverage $coverage% is below 80%"
            exit 1
          fi
```

**Rationale**: Separate lint and test jobs enable parallel execution. Using `stable` Go version ensures latest features. Shell script for coverage check is simple and reliable.

**Alternatives Considered**:
- codecov.io integration → Rejected; adds external dependency and complexity for MVP
- Matrix builds (multiple Go versions) → Deferred; single stable version sufficient for MVP
- Separate coverage job → Rejected; integrated with test job is simpler

---

### 8. Makefile Targets

**Question**: What Make targets should we provide for local development?

**Research**: Reviewed common Go project Makefile patterns and developer workflows.

**Decision**: Essential targets:
```makefile
.PHONY: test lint build coverage clean

test:
	go test -v -race ./...

lint:
	golangci-lint run

build:
	go build -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$coverage%"; \
	if [ $$(echo "$$coverage < 80" | bc -l) -eq 1 ]; then \
		echo "ERROR: Coverage $$coverage% is below 80%"; \
		exit 1; \
	fi

clean:
	rm -f coverage.out coverage.html
```

**Rationale**: Mirrors CI workflow for consistency. Coverage target generates HTML for visual inspection. All targets use standard Go tooling.

**Alternatives Considered**:
- `make all` → Rejected; unclear what "all" means for library (no binary to build)
- `make install` → Deferred; users install via `go get`, not make
- `make fmt` → Rejected; developers should use editor integration, CI enforces gofmt

---

## Technology Stack Summary

| Component | Choice | Version |
|-----------|--------|---------|
| **Language** | Go | Latest stable (1.21+) |
| **HTTP Client** | `net/http` (stdlib) | Stdlib |
| **Testing** | `testing` (stdlib) + `httptest` | Stdlib |
| **Linting** | golangci-lint | Latest |
| **CI/CD** | GitHub Actions | N/A |
| **Build** | make + go toolchain | N/A |
| **Dependencies** | None (stdlib only) | N/A |

## Dependencies Justification

**Zero external dependencies**: The SDK uses only Go standard library. This aligns with constitution principle of vetting dependencies and ensures:
- No supply chain security risks
- Minimal maintenance burden
- Fast compilation
- Easy auditing

All required functionality (HTTP, JSON, testing) is available in stdlib.

## Performance Considerations

**Connection Pooling**: Reusing `http.Client` enables automatic connection pooling (default 100 idle connections per host).

**Memory Allocation**: Weather response structures use value types (not pointers) where possible to reduce GC pressure.

**Concurrency**: Semaphore pattern with buffered channel is lock-free and scales to thousands of goroutines.

**Benchmarking**: Defer until after MVP; no performance issues expected for I/O-bound HTTP requests.

## Security Considerations

**Input Validation**: All coordinates validated before network requests to prevent injection attacks.

**HTTPS Only**: Base URL hardcoded to `https://` to prevent downgrade attacks.

**Timeout Protection**: Default 10s timeout prevents resource exhaustion from slow/malicious servers.

**No Credential Storage**: API requires no authentication; no secrets to manage.

## Documentation Plan

**godoc**: All exported types, functions, and methods documented with examples.

**README.md**: Installation, quick start, usage examples, contributing guidelines.

**examples/**: Working code samples for common use cases.

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Open Meteo API changes breaking contract | Low | High | Integration tests detect breakage; version SDK with SemVer |
| Rate limiting on free tier | Medium | Medium | Document limits; no SDK-level mitigation per clarifications |
| Missing weather parameters | High | Low | Use zero values and document per clarifications |
| Coordinate validation edge cases | Low | Low | Comprehensive unit tests for boundaries |

### 4. Toolchain Versioning

**Question**: What are the latest stable versions of the Go language and GitHub Actions components as of December 2025?

**Research**: Performed web search for latest releases of Go, actions/checkout, actions/setup-go, and golangci-lint-action.

**Decision**:
- **Go Language**: 1.25.5 (Latest Stable)
- **actions/checkout**: v6.0.1
- **actions/setup-go**: v6.1.0
- **golangci-lint-action**: v9.2.0

**Rationale**: Using the latest stable versions ensures access to the newest language features, performance improvements, and security patches. It also ensures compatibility with modern CI environments.

**Alternatives Considered**:
- Using older LTS versions -> Rejected to maximize longevity of the codebase without immediate upgrade needs.

## Next Steps (Phase 1)

With research complete, proceed to:
1. Define data model in `data-model.md`
2. Design API contracts in `contracts/`
3. Create quick start guide in `quickstart.md`
4. Update agent context with technology decisions
