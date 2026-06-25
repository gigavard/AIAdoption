# ExTodoGolang - Spec-Driven Todo App in Go

Go Todo App following TDD and spec-driven development methodology. Built as part of the AIAdoption educational exercise.

## Quick Start

```bash
# Install dependencies
make deps

# Run tests (TDD-first)
make test

# Build and run
make build
make run

# Or run directly
make run
```

## Project Structure

```
ExTodoGolang/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management (SPEC-001)
│   ├── domain/         # Domain models (SPEC-002)
│   ├── storage/        # SQLite persistence (SPEC-002)
│   ├── http/           # HTTP routes & handlers (SPEC-003)
│   ├── observability/  # Logging, tracing, metrics (SPEC-006)
│   └── errors/         # RFC 7807 error handling (SPEC-004)
├── pkg/
│   ├── logger/         # JSON logging
│   └── errors/         # Error types
├── web/
│   ├── public/         # Static JS/CSS (vanilla + Tailwind)
│   └── templates/      # HTML templates
├── migrations/         # SQLite schema migrations
├── tests/
│   ├── unit/          # Unit test suite
│   └── integration/   # Integration test suite
└── Makefile           # Build and test targets
```

## Development Workflow

1. **SPEC-001**: Project bootstrap ✅ (you are here)
2. **SPEC-002**: Todo domain model + SQLite persistence
3. **SPEC-003**: HTTP CRUD API + OpenAPI spec
4. **SPEC-004**: RFC 7807 error handling
5. **SPEC-005**: Vanilla JS + Tailwind frontend
6. **SPEC-006**: Observability (logging, telemetry, Prometheus)
7. **SPEC-007**: Graceful shutdown
8. **SPEC-008**: TDD strategy & test pyramid

## References

- **Documentation**: See `https://github.com/gigavard/AIKnow` for round table notes and exercise definitions
- **Spec Tracking**: All specs tracked in DocMind project `Ex_todo_golang`
- **Round Table Source**: From "2026-02-17 Round Table" Esercizio 3

## Compliance

- ✅ Go 1.23+
- ✅ chi router (will add in SPEC-003)
- ✅ SQLite with sqlx (will add in SPEC-002)
- ✅ Vanilla JS + Tailwind CSS
- ✅ OpenAPI 3.0 specification
- ✅ RFC 7807 problem details for errors
- ✅ OpenTelemetry instrumentation
- ✅ Prometheus metrics
- ✅ JSON structured logging
- ✅ Graceful HTTP shutdown
