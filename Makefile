.PHONY: help fmt lint test test-coverage build clean fuzz

# Default target
help:
	@echo "Available targets:"
	@echo "  make fmt            - Format Go code"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make fuzz           - Run all fuzz tests (or specific test with FUZZ_TEST=FuzzName)"
	@echo "  make build          - Build the project"
	@echo "  make clean          - Clean build artifacts"

fmt:
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "Done!"

lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Done!"

test:
	@echo "Running tests..."
ifeq ($(OS),Windows_NT)
	@go test -v -race ./...
else
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
endif
	@echo "Done!"

test-coverage:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@echo "Done! Run 'go tool cover -html=coverage.out' to view detailed coverage."

build:
	@echo "Building project..."
	@go build ./...
	@echo "Done!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -f coverage.out
	@go clean ./...
	@echo "Done!"

fuzz:
ifdef FUZZ_TEST
	@echo "Running fuzz test: $(FUZZ_TEST)..."
	@go test -run=^$$ -fuzz=$(FUZZ_TEST) -fuzztime=30s .
	@echo "Done!"
else
	@echo "Running all fuzz tests..."
	@for test in $$(grep -h '^func Fuzz' *_test.go | sed 's/func \(Fuzz[^(]*\).*/\1/'); do \
		echo "Running $$test..."; \
		go test -run=^$$ -fuzz=$$test -fuzztime=30s . || exit 1; \
	done
	@echo "Done!"
endif
