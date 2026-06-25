.PHONY: help build run test clean dev lint

help:
	@echo "ExTodoGolang - Go Todo App"
	@echo ""
	@echo "Targets:"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  test      - Run all tests"
	@echo "  test-unit - Run unit tests only"
	@echo "  test-int  - Run integration tests only"
	@echo "  dev       - Run with auto-reload (requires air)"
	@echo "  lint      - Run linter (requires golangci-lint)"
	@echo "  clean     - Remove build artifacts"

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

test:
	go test -v -race -coverprofile=coverage.out ./...

test-unit:
	go test -v -short ./...

test-int:
	go test -v -run Integration ./...

dev:
	air

lint:
	golangci-lint run

clean:
	rm -f bin/server coverage.out

deps:
	go mod download
	go mod tidy
