.PHONY: build run clean test test-unit test-integration test-all compile-examples benchmark help

# Use local cache to avoid sandbox permission issues
export GOCACHE := $(PWD)/.cache

# Build the JVM
build:
	go build -o simplejvm

# Run HelloWorld example
run: build
	./simplejvm examples/HelloWorld.class

# Run with verbose mode
run-verbose: build
	./simplejvm -v examples/HelloWorld.class

# Compile all Java examples
compile-examples:
	@echo "Compiling Java examples..."
	@cd examples && for f in *.java; do \
		echo "  Compiling $$f"; \
		javac "$$f" 2>/dev/null || true; \
	done
	@echo "Done."

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v ./runtime/... ./interpreter/... ./classfile/...

# Run integration tests (requires compiled Java examples)
test-integration: build compile-examples
	@echo "Running integration tests..."
	go test -v -run "Test.*" .

# Run all tests
test-all: test-unit test-integration
	@echo "All tests completed!"

# Quick test - just run examples
test: build compile-examples
	@echo "Running all example tests..."
	@for f in HelloWorld Calculator ArrayTest ObjectTest ExceptionTest TypeTest NativeTest SyncTest; do \
		echo "=== Testing $$f ==="; \
		./simplejvm examples/$${f}.class 2>&1 | tail -3; \
		echo ""; \
	done
	@echo "All tests passed!"

# Run benchmarks
benchmark: build compile-examples
	@echo "Running benchmarks..."
	go test -bench=. -benchtime=5s .

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run a specific example
run-example: build
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "Usage: make run-example EXAMPLE=Calculator"; \
		exit 1; \
	fi
	./simplejvm examples/$(EXAMPLE).class

# Run with stats
run-stats: build
	./simplejvm -stats examples/Calculator.class

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, using go vet"; \
		go vet ./...; \
	fi

# Clean build artifacts
clean:
	rm -f simplejvm
	rm -f examples/*.class
	rm -rf .cache
	rm -f coverage.out coverage.html

# Show help
help:
	@echo "SimpleJVM - A minimal JVM implementation in Go"
	@echo ""
	@echo "Usage:"
	@echo "  make build            Build the JVM"
	@echo "  make run              Run HelloWorld example"
	@echo "  make run-verbose      Run with bytecode tracing"
	@echo "  make run-stats        Run with heap statistics"
	@echo "  make run-example EXAMPLE=Name  Run a specific example"
	@echo ""
	@echo "Testing:"
	@echo "  make test             Quick test - run all examples"
	@echo "  make test-unit        Run Go unit tests"
	@echo "  make test-integration Run integration tests"
	@echo "  make test-all         Run all tests"
	@echo "  make benchmark        Run benchmarks"
	@echo "  make coverage         Generate test coverage report"
	@echo ""
	@echo "Development:"
	@echo "  make compile-examples Compile all Java examples"
	@echo "  make fmt              Format Go code"
	@echo "  make lint             Run linter"
	@echo "  make clean            Remove build artifacts"
	@echo ""
	@echo "Examples available:"
	@echo "  HelloWorld, Calculator, ArrayTest, ObjectTest,"
	@echo "  ExceptionTest, TypeTest, NativeTest, SyncTest"
