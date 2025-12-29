# Tasks: Open Meteo Current Weather SDK

**Branch**: `001-current-weather-sdk` | **Date**: 2025-12-29
**Input**: Design documents from `/specs/001-current-weather-sdk/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: This feature does NOT explicitly request TDD or test-first development. Test tasks are included but can be executed alongside or after implementation per standard Go testing practices.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Initialize Go module at /workspaces/open-meteo-weather-sdk/go.mod
- [X] T002 [P] Create .gitignore file at /workspaces/open-meteo-weather-sdk/.gitignore with Go patterns
- [X] T003 [P] Create Makefile at /workspaces/open-meteo-weather-sdk/Makefile with test, lint, coverage, clean targets
- [X] T004 [P] Create GitHub Actions workflow at /workspaces/open-meteo-weather-sdk/.github/workflows/ci.yml using Go 1.25.5, actions/checkout@v6.0.1, actions/setup-go@v6.1.0, golangci-lint-action@v9.2.0
- [X] T005 [P] Create initial README.md at /workspaces/open-meteo-weather-sdk/README.md with project description and placeholder for usage

**Checkpoint**: Project structure ready - can begin implementation work

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T006 Define CurrentWeather struct with all 15 weather parameter fields in /workspaces/open-meteo-weather-sdk/weather.go
- [X] T007 [P] Define internal weatherResponse struct for API JSON unmarshaling in /workspaces/open-meteo-weather-sdk/weather.go
- [X] T008 [P] Define Client struct with httpClient, baseURL, semaphore fields in /workspaces/open-meteo-weather-sdk/client.go
- [X] T009 [P] Define ErrorType constants (Validation, Network, API) in /workspaces/open-meteo-weather-sdk/errors.go
- [X] T010 [P] Define Error struct with Type, Message, Cause fields in /workspaces/open-meteo-weather-sdk/errors.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Fetch Current Weather by Coordinates (Priority: P1) 🎯 MVP

**Goal**: Developers can retrieve current weather data for a specific location using latitude and longitude with a single method call

**Independent Test**: Call `GetCurrentWeather(ctx, 52.52, 13.41)` with valid coordinates and verify all 15 weather parameters are returned with reasonable values

### Implementation for User Story 1

- [X] T011 [P] [US1] Implement NewClient constructor with default configuration (10s timeout, base URL, 10-capacity semaphore) in /workspaces/open-meteo-weather-sdk/client.go
- [X] T012 [P] [US1] Implement GetCurrentWeather method signature with context, lat, lon parameters in /workspaces/open-meteo-weather-sdk/client.go
- [X] T013 [US1] Implement coordinate validation logic in GetCurrentWeather (lat: -90 to 90, lon: -180 to 180) in /workspaces/open-meteo-weather-sdk/client.go
- [X] T014 [US1] Implement HTTP request construction with query parameters (latitude, longitude, current_weather=true, metric units) in /workspaces/open-meteo-weather-sdk/client.go
- [X] T015 [US1] Implement HTTP request execution with context and timeout in /workspaces/open-meteo-weather-sdk/client.go
- [X] T016 [US1] Implement JSON response parsing into CurrentWeather struct in /workspaces/open-meteo-weather-sdk/client.go
- [X] T017 [US1] Implement semaphore-based concurrency control (max 10 requests) in GetCurrentWeather in /workspaces/open-meteo-weather-sdk/client.go
- [X] T018 [P] [US1] Add godoc comments to Client, NewClient, and GetCurrentWeather in /workspaces/open-meteo-weather-sdk/client.go
- [X] T019 [P] [US1] Add godoc comments to CurrentWeather and all fields in /workspaces/open-meteo-weather-sdk/weather.go

### Tests for User Story 1

- [X] T020 [P] [US1] Create weather_test.go with test for JSON unmarshaling of complete API response in /workspaces/open-meteo-weather-sdk/weather_test.go
- [X] T021 [P] [US1] Add test for JSON unmarshaling with null/missing weather parameters (zero values) in /workspaces/open-meteo-weather-sdk/weather_test.go
- [X] T022 [US1] Create client_test.go with httptest mock server for successful GetCurrentWeather call in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T023 [US1] Add test for GetCurrentWeather with valid boundary coordinates (lat=90, lon=180) in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T024 [US1] Add test for concurrent GetCurrentWeather calls (up to 10 simultaneous) in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T025 [US1] Add test for context cancellation during GetCurrentWeather in /workspaces/open-meteo-weather-sdk/client_test.go

**Checkpoint**: User Story 1 complete - SDK can fetch weather data successfully

---

## Phase 4: User Story 2 - Handle Errors Gracefully (Priority: P2)

**Goal**: Provide clear, actionable error messages distinguishing between validation errors, network errors, and API errors

**Independent Test**: Trigger each error condition (invalid coordinates, network failure, API 500) and verify appropriate error type and message are returned

### Implementation for User Story 2

- [X] T026 [P] [US2] Implement Error.Error() method returning formatted message in /workspaces/open-meteo-weather-sdk/errors.go
- [X] T027 [P] [US2] Implement Error.Unwrap() method for error chain inspection in /workspaces/open-meteo-weather-sdk/errors.go
- [X] T028 [P] [US2] Add godoc comments to ErrorType, Error, and methods in /workspaces/open-meteo-weather-sdk/errors.go
- [X] T029 [US2] Update GetCurrentWeather to return ErrorTypeValidation for invalid coordinates in /workspaces/open-meteo-weather-sdk/client.go
- [X] T030 [US2] Update GetCurrentWeather to return ErrorTypeValidation for concurrent request limit exceeded in /workspaces/open-meteo-weather-sdk/client.go
- [X] T031 [US2] Update GetCurrentWeather to return ErrorTypeNetwork for network/HTTP failures in /workspaces/open-meteo-weather-sdk/client.go
- [X] T032 [US2] Update GetCurrentWeather to return ErrorTypeAPI for non-200 HTTP status codes in /workspaces/open-meteo-weather-sdk/client.go
- [X] T033 [US2] Update GetCurrentWeather to return ErrorTypeAPI for JSON parsing failures in /workspaces/open-meteo-weather-sdk/client.go

### Tests for User Story 2

- [X] T034 [P] [US2] Create errors_test.go with test for Error.Error() formatting in /workspaces/open-meteo-weather-sdk/errors_test.go
- [X] T035 [P] [US2] Add test for Error.Unwrap() with wrapped errors in /workspaces/open-meteo-weather-sdk/errors_test.go
- [X] T036 [P] [US2] Add test for errors.Is and errors.As compatibility in /workspaces/open-meteo-weather-sdk/errors_test.go
- [X] T037 [P] [US2] Add test in client_test.go for ErrorTypeValidation with invalid latitude (999.0) in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T038 [P] [US2] Add test in client_test.go for ErrorTypeValidation with invalid longitude (-999.0) in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T039 [P] [US2] Add test in client_test.go for ErrorTypeValidation with concurrent limit exceeded in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T040 [P] [US2] Add test in client_test.go for ErrorTypeNetwork with httptest server timeout in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T041 [P] [US2] Add test in client_test.go for ErrorTypeAPI with HTTP 400 response in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T042 [P] [US2] Add test in client_test.go for ErrorTypeAPI with HTTP 500 response in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T043 [P] [US2] Add test in client_test.go for ErrorTypeAPI with malformed JSON response in /workspaces/open-meteo-weather-sdk/client_test.go

**Checkpoint**: User Story 2 complete - SDK provides typed, actionable errors

---

## Phase 5: User Story 3 - Configure Request Options (Priority: P3)

**Goal**: Allow developers to customize SDK behavior (timeout, HTTP client, base URL) via functional options

**Independent Test**: Create SDK clients with different configurations and verify requests honor those settings

### Implementation for User Story 3

- [X] T044 [P] [US3] Define Option function type in /workspaces/open-meteo-weather-sdk/options.go
- [X] T045 [P] [US3] Implement WithTimeout option in /workspaces/open-meteo-weather-sdk/options.go
- [X] T046 [P] [US3] Implement WithHTTPClient option in /workspaces/open-meteo-weather-sdk/options.go
- [X] T047 [P] [US3] Implement WithBaseURL option in /workspaces/open-meteo-weather-sdk/options.go
- [X] T048 [P] [US3] Add godoc comments to Option type and all option functions in /workspaces/open-meteo-weather-sdk/options.go
- [X] T049 [US3] Update NewClient to accept variadic Option parameters and apply them in /workspaces/open-meteo-weather-sdk/client.go
- [X] T050 [US3] Update Client initialization to support option-based configuration in /workspaces/open-meteo-weather-sdk/client.go

### Tests for User Story 3

- [X] T051 [P] [US3] Create options_test.go with test for WithTimeout option in /workspaces/open-meteo-weather-sdk/options_test.go
- [X] T052 [P] [US3] Add test in options_test.go for WithHTTPClient option in /workspaces/open-meteo-weather-sdk/options_test.go
- [X] T053 [P] [US3] Add test in options_test.go for WithBaseURL option in /workspaces/open-meteo-weather-sdk/options_test.go
- [X] T054 [P] [US3] Add test in options_test.go for multiple options combined in /workspaces/open-meteo-weather-sdk/options_test.go
- [X] T055 [P] [US3] Add test in client_test.go for custom timeout enforcement (request exceeds timeout) in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T056 [P] [US3] Add test in client_test.go for custom HTTP client with transport in /workspaces/open-meteo-weather-sdk/client_test.go
- [X] T057 [P] [US3] Add test in client_test.go for custom base URL in /workspaces/open-meteo-weather-sdk/client_test.go

**Checkpoint**: All user stories complete - SDK is fully configurable

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, examples, and final validation

- [X] T058 [P] Create examples/basic/main.go with basic usage example
- [X] T059 Update README.md with installation instructions, usage examples, and API reference in /workspaces/open-meteo-weather-sdk/README.md
- [X] T060 [P] Verify all exported symbols have godoc comments across all files
- [X] T061 Run make lint and fix any linting issues
- [X] T062 Run make coverage and verify 80% minimum coverage
- [X] T063 Validate quickstart.md examples can be executed successfully
- [X] T064 Create LICENSE file at /workspaces/open-meteo-weather-sdk/LICENSE

**Checkpoint**: SDK ready for v0.1.0 release

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion (go.mod must exist) - BLOCKS all user stories
- **User Stories (Phases 3-5)**: All depend on Foundational phase completion
  - US1 (P1) can start immediately after Phase 2
  - US2 (P2) can start after Phase 2, but benefits from US1 being complete for integration
  - US3 (P3) can start after Phase 2, requires US1 complete for full testing
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories - THIS IS THE MVP
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Enhances US1 but is independently testable with mocked errors
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Extends US1 but can be tested independently with configuration

### Within Each User Story

- **US1**: 
  - Implementation tasks T011-T019 can proceed mostly in parallel (T013-T017 depend on T012)
  - Tests T020-T025 can run in parallel once implementation complete
- **US2**:
  - Error implementation tasks T026-T028 are parallel
  - Client updates T029-T033 depend on T026-T028
  - All tests T034-T043 can run in parallel
- **US3**:
  - Options implementation T044-T048 are fully parallel
  - Client updates T049-T050 depend on T044-T047
  - All tests T051-T057 can run in parallel

### Parallel Opportunities

**Within Phase 1 (Setup)**: Tasks T002, T003, T004, T005 can all run in parallel after T001

**Within Phase 2 (Foundational)**: Tasks T007, T008, T009, T010 can run in parallel after T006

**Within Phase 3 (US1)**: 
- Tasks T011, T012 can start together
- Tasks T018, T019 can run in parallel
- All test tasks T020-T025 can run in parallel

**Within Phase 4 (US2)**:
- Tasks T026, T027, T028 fully parallel
- All test tasks T034-T043 fully parallel

**Within Phase 5 (US3)**:
- Tasks T044, T045, T046, T047, T048 fully parallel
- All test tasks T051-T057 fully parallel

**Within Phase 6 (Polish)**:
- Tasks T058, T060, T063, T064 can run in parallel

**Cross-Story Parallelism**: If multiple developers available, US2 and US3 can proceed in parallel once US1 completes

---

## Parallel Example: User Story 1 Implementation

```bash
# After Phase 2 complete, these can run simultaneously:

# Developer A: Core client implementation
# T011-T017 (in client.go)

# Developer B: Documentation
# T018-T019 (godoc comments)

# Developer C: Tests
# T020-T025 (once implementation exists)
```

---

## Implementation Strategy

**MVP Scope**: Phase 1 + Phase 2 + Phase 3 (User Story 1 only)
- This delivers basic weather fetching capability
- Estimated: ~600 LOC implementation + ~600 LOC tests
- Timeline: 1-2 days for single developer

**Incremental Delivery**:
1. MVP (US1): Basic weather fetching
2. +US2: Production-ready error handling
3. +US3: Full configurability

**Suggested MVP**: Complete through Phase 3 (User Story 1), then validate with quickstart scenarios before proceeding to P2/P3.

---

## Task Summary

- **Total Tasks**: 64
- **Setup**: 5 tasks (1 sequential + 4 parallel)
- **Foundational**: 5 tasks (1 + 4 parallel)
- **User Story 1 (P1)**: 15 tasks (9 implementation + 6 tests)
- **User Story 2 (P2)**: 18 tasks (8 implementation + 10 tests)
- **User Story 3 (P3)**: 14 tasks (7 implementation + 7 tests)
- **Polish**: 7 tasks (5 parallel + 2 sequential)

**Parallel Opportunities**: 42 tasks marked [P] can run in parallel within their phases

**Independent Test Validation**:
- US1: Fully testable with basic GetCurrentWeather calls
- US2: Fully testable with error injection scenarios
- US3: Fully testable with configuration options

**Format Validation**: ✅ All tasks follow checklist format with:
- Checkbox `- [ ]`
- Task ID (T001-T064)
- [P] marker where appropriate
- [US1], [US2], [US3] labels for user story phases
- Exact file paths included in descriptions
