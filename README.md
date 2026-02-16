# gogem

A collection of reusable Go packages for building web applications. Each package is an independent Go module, versioned and released separately.

## Installation

Install only the packages you need:

```bash
go get github.com/bernardinorafael/gogem/pkg/fault@latest
go get github.com/bernardinorafael/gogem/pkg/httputil@latest
go get github.com/bernardinorafael/gogem/pkg/uid@latest
```

## Packages

| Package | Description | Dependencies |
|---------|-------------|-------------|
| [`fault`](./pkg/fault) | Standardized REST error type with HTTP codes, tags, and field-level validation errors | - |
| [`httputil`](./pkg/httputil) | Request parsing, JSON response writing, and `WithValidation[T]` generic middleware | fault |
| [`pagination`](./pkg/pagination) | Generic `Paginated[T]` container with computed metadata | - |
| [`uid`](./pkg/uid) | K-Sortable unique identifier generation with optional prefixes | - |
| [`function`](./pkg/function) | Generic `Map` and `ForEach` utilities | - |
| [`apiutil`](./pkg/apiutil) | Generic `Expandable[T]` for API responses (marshals as ID or full object) | uid |
| [`dbutil`](./pkg/dbutil) | PostgreSQL helpers: constraint violation detection, JSONB type, transaction wrapper | fault, sqlx, lib/pq |
| [`cache`](./pkg/cache) | Redis wrapper with generic `GetOrSet[T]` pattern | fault, go-redis, charmbracelet/log |
| [`crypto`](./pkg/crypto) | Password hashing (bcrypt), JWT tokens (HS256), OTP generation (HMAC-SHA256) | uid, golang-jwt, x/crypto |
| [`queue`](./pkg/queue) | AWS SQS client wrapper (publish, consume, delete) | aws-sdk-go-v2 |
| [`logger`](./pkg/logger) | Structured logging (JSON in production, text in development) | charmbracelet/log |
| [`server`](./pkg/server) | HTTP server with chi router, sensible defaults, and graceful shutdown | go-chi |

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
make help      # show all available commands
make vet       # run go vet on all modules
make test      # run tests on all modules
make tidy      # run go mod tidy on all modules
make fmt       # format all Go files
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
make tags                    # list latest tag for each module
make changed                 # show modules changed since last tag
make tag m=fault b=patch     # pkg/fault/v0.1.0 -> pkg/fault/v0.1.1
make tag m=fault b=minor     # pkg/fault/v0.1.0 -> pkg/fault/v0.2.0
make tag m=fault b=major     # pkg/fault/v0.1.0 -> pkg/fault/v1.0.0
```

For modules with internal dependencies (Layer 1), tag the dependencies first:

```bash
make tag m=uid b=minor       # tag dependency first
make tag m=crypto b=minor    # then tag dependent module
```

Consumers install specific modules at specific versions:

```bash
go get github.com/bernardinorafael/gogem/pkg/fault@v0.1.0
```

## License

MIT
