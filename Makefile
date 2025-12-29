.PHONY: test lint coverage clean help

# Run tests with race detector
test:
	@echo "Running tests..."
	go test -v -race ./...

# Run linter (requires golangci-lint installed)
lint:
	@echo "Running linter..."
	golangci-lint run

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out.tmp ./...
	@grep -v "/examples/" coverage.out.tmp > coverage.out || true
	@rm -f coverage.out.tmp
	@go tool cover -html=coverage.out -o coverage.html
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$coverage%"; \
	if [ "$$(echo "$$coverage" | awk '{if ($$1 < 80) print 1; else print 0}')" -eq 1 ]; then \
		echo "ERROR: Coverage $$coverage% is below 80% threshold"; \
		exit 1; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	go clean

# Display available targets
help:
	@echo "Available targets:"
	@echo "  test     - Run tests with race detector"
	@echo "  lint     - Run golangci-lint"
	@echo "  coverage - Generate coverage report (requires 80%)"
	@echo "  clean    - Remove build artifacts"
	@echo "  help     - Display this help message"
