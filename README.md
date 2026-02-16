# gogem

A collection of reusable Go packages for building web applications. Each package is an independent Go module, versioned and released separately.

## Installation

Install only the packages you need:

```bash
go get github.com/bernardinorafael/gogem/fault@latest
go get github.com/bernardinorafael/gogem/httputil@latest
go get github.com/bernardinorafael/gogem/uid@latest
```

## Packages

| Package | Description | Dependencies |
|---------|-------------|-------------|
| [`fault`](./fault) | Standardized REST error type with HTTP codes, tags, and field-level validation errors | - |
| [`httputil`](./httputil) | Request parsing, JSON response writing, and `WithValidation[T]` generic middleware | fault |
| [`pagination`](./pagination) | Generic `Paginated[T]` container with computed metadata | - |
| [`uid`](./uid) | K-Sortable unique identifier generation with optional prefixes | - |
| [`function`](./function) | Generic `Map` and `ForEach` utilities | - |
| [`apiutil`](./apiutil) | Generic `Expandable[T]` for API responses (marshals as ID or full object) | uid |
| [`dbutil`](./dbutil) | PostgreSQL helpers: constraint violation detection, JSONB type, transaction wrapper | fault, sqlx, lib/pq |
| [`cache`](./cache) | Redis wrapper with generic `GetOrSet[T]` pattern | fault, go-redis, charmbracelet/log |
| [`crypto`](./crypto) | Password hashing (bcrypt), JWT tokens (HS256), OTP generation (HMAC-SHA256) | uid, golang-jwt, x/crypto |
| [`queue`](./queue) | AWS SQS client wrapper (publish, consume, delete) | aws-sdk-go-v2 |
| [`logger`](./logger) | Structured logging (JSON in production, text in development) | charmbracelet/log |
| [`server`](./server) | HTTP server with chi router, sensible defaults, and graceful shutdown | go-chi |

## Development

This is a multi-module repository. Each package has its own `go.mod` and is developed using [Go workspaces](https://go.dev/doc/tutorial/workspaces).

### Prerequisites

- Go 1.24.1+

### Setup

Clone the repository. The `go.work` file at the root links all modules for local development:

```bash
git clone https://github.com/bernardinorafael/gogem.git
cd gogem
```

All modules are already listed in `go.work`, so cross-module imports resolve locally.

### Commands

```bash
# Vet a specific module
cd fault && go vet ./...

# Run tests for a specific module
cd fault && go test ./...

# Tidy a specific module's dependencies
cd fault && go mod tidy

# Format all code
gofmt -w .
```

### Dependency Graph

```
Layer 0 (no internal deps):
  fault, pagination, uid, function, queue, logger, server

Layer 1 (depends on Layer 0):
  httputil → fault
  dbutil   → fault
  cache    → fault
  apiutil  → uid
  crypto   → uid
```

### Releasing a Module

Each module is versioned independently using git tags with module prefix:

```bash
# Tag a module
git tag fault/v0.1.0
git push origin fault/v0.1.0

# For modules with internal deps (Layer 1), tag Layer 0 deps first
git tag uid/v0.1.0
git push origin uid/v0.1.0

# Then update go.mod to point to the published version and tag
git tag crypto/v0.1.0
git push origin crypto/v0.1.0
```

Consumers install specific modules at specific versions:

```bash
go get github.com/bernardinorafael/gogem/fault@v0.1.0
```

## License

MIT
