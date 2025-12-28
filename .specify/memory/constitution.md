<!--
Sync Impact Report:
- Version Change: [NEW] → 1.0.0
- Principles Added: 6 core principles established
  1. Code Quality (Effective Go)
  2. Testing Standards
  3. User Experience Consistency
  4. Performance Requirements
  5. Documentation Standards
  6. Release & Build Standards
- Sections Added: Technical Stack, Development Workflow, Governance
- Templates Status:
  ✅ plan-template.md - reviewed, constitution checks align with principles
  ✅ spec-template.md - reviewed, requirements align with principles
  ✅ tasks-template.md - reviewed, task organization aligns with principles
  ✅ agent-file-template.md - reviewed, structure supports governance
- Follow-up TODOs: None - all placeholders filled
- Rationale: Initial constitution establishment (MINOR version 1.0.0)
-->

# Open-Meteo Weather SDK Constitution

## Core Principles

### I. Code Quality (Effective Go)

We strictly adhere to the idioms and best practices outlined in [Effective Go](https://go.dev/doc/effective_go). Code MUST be formatted with `gofmt`. Naming conventions, error handling, and concurrency patterns MUST follow Go community standards. Clarity and simplicity are preferred over cleverness. Public APIs MUST be documented.

**Rationale**: Consistency across the Go ecosystem improves maintainability, reduces onboarding friction, and ensures code reviews focus on logic rather than style debates.

### II. Testing Standards

Testing is mandatory. All packages MUST have unit tests (`_test.go`) colocated with source code. Test coverage MUST meet a minimum of 80% for unit tests. Integration tests are REQUIRED for external interactions (APIs, file systems). Tests MUST be deterministic and fast.

**Rationale**: High test coverage catches regressions early, enables confident refactoring, and serves as executable documentation. Colocated tests reduce cognitive overhead during development.

### III. User Experience Consistency

The CLI and output MUST be consistent and predictable. Use standard flags and arguments. Output SHOULD be human-readable by default, with options for machine-readable formats (e.g., JSON) where appropriate. Error messages MUST be actionable and clear to the end-user, distinguishing between user errors and system failures.

**Rationale**: Consistent interfaces reduce learning curves and prevent user frustration. Actionable error messages reduce support burden and improve troubleshooting efficiency.

### IV. Performance Requirements

The application MUST be efficient with resources (CPU, Memory). Avoid unnecessary allocations in hot paths. Network operations MUST have timeouts. Performance critical paths SHOULD be benchmarked. Latency for user-facing operations SHOULD be minimized.

**Rationale**: Weather data applications are often used in resource-constrained environments or high-frequency scenarios. Efficient resource usage ensures scalability and responsiveness.

### V. Documentation Standards

A `README.md` file is REQUIRED at the root of the repository. It MUST document the current state of the project, including what it does, why it is useful, how to get started, and how to contribute. This file MUST be updated after implementation is complete and before a Pull Request is created to reflect any changes in functionality, usage, or configuration.

**Rationale**: Documentation is the first touchpoint for users and contributors. Keeping it current prevents confusion and reduces repetitive support questions.

### VI. Release & Build Standards

Builds MUST be reproducible, secure, and automated. Releases MUST be triggered by tags following Semantic Versioning (SemVer). The release pipeline MUST enforce linting and testing before publishing. Artifacts MUST include version metadata (embedded at build time) and be cross-compiled for all supported target platforms. Integrity checks (e.g., checksums) MUST be provided for all published artifacts.

**Rationale**: Reproducible builds enable auditability and security verification. Automated pipelines prevent human error and ensure consistent quality gates across all releases.

## Technical Stack

**Language**: Go (Latest Stable)

**Dependency Management**: Go Modules

**Linter**: `golangci-lint` (Standard configuration)

**Build Tool**: Standard `go build`

## Development Workflow

**Branching**: Feature branches off `main`.

**Commits**: Follow Conventional Commits specification.

**Review**: Pull Request required for all changes. Code review MUST verify compliance with Core Principles.

**CI/CD**: Automated tests and linters MUST pass before merging.

## Governance

This Constitution supersedes all other practices. Amendments require documentation, approval, and a migration plan.

**Rules**:

1. All PRs/reviews MUST verify compliance with "Effective Go" and this Constitution.
2. Complexity MUST be justified.
3. New dependencies MUST be vetted for license and maintenance status.
4. Versioning follows Semantic Versioning (SemVer).

**Version**: 1.0.0 | **Ratified**: 2025-12-28 | **Last Amended**: 2025-12-28
