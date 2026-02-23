# mem-local

A small Go project providing a simple in-memory local store. This repository contains a minimal example application and a lightweight `store` implementation intended for local, ephemeral usage and testing.

**Status:** minimal example / reference implementation

**Prerequisites:**

**Quick Start**

Build and run the example locally:

```bash
go build -o mem-local .
./mem-local
# or, during development:
go run main.go
```

Check `main.go` for how the store is constructed and used at runtime.

**Project Layout**

**Usage Notes**

**Contributing**

**License**

# mem-local

An in-memory key-value store for Go, designed for local, ephemeral usage, testing, and prototyping. Includes a minimal example application and a thread-safe store implementation.

**Status:** Minimal reference implementation

**Prerequisites:**
- Go 1.20+ (tested with Go 1.23)

## Quick Start

Build and run the example:

```bash
go build -o mem-local .
./mem-local
# or, during development:
go run example.go
```

See [example.go](example.go) for usage.

## API Overview

The store provides the following methods:

- `Init() error` — Initialize the store (clears all data)
- `Set(key, value string) error` — Store a value
- `Get(key string) (string, bool)` — Retrieve a value
- `Delete(key string) error` — Remove a key
- `Exists(key string) bool` — Check if a key exists

All methods are safe for concurrent use.

## Example Usage

```go
import "github.com/v-72/mem-local/store"

s := &store.MemLocalStore{}
_ = s.Init()
s.Set("foo", "bar")
val, ok := s.Get("foo")
exists := s.Exists("foo")
s.Delete("foo")
```

## Project Layout
- [example.go](example.go): Example program / entrypoint
- [store/mem_local.go](store/mem_local.go): In-memory store implementation
- [store/mem_local_test.go](store/mem_local_test.go): Tests
- [go.mod](go.mod): Go module definition

## Usage Notes
- Data is not persisted across restarts.
- Designed for testing, prototypes, or local caches where persistence is unnecessary.

## Contributing
- Issues and PRs are welcome. Please keep changes small and focused.

## License
- No license included. Add one (`LICENSE`) if you intend to open-source this project.
