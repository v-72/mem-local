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
- `Set(key, value string, ttl *int) error` — Store a value with optional TTL (milliseconds)
- `Get(key string) (string, bool)` — Retrieve a value
- `Delete(key string) error` — Remove a key
- `Exists(key string) bool` — Check if a key exists
- `SetMany(kv map[string]string, ttl *int) error` — Store multiple values
- `GetMany(keys []string) []string` — Retrieve multiple values

All methods are safe for concurrent use.

## Example Usage

### Basic Operations

```go
import "github.com/v-72/mem-local/store"

s := &store.MemLocalStore{}
_ = s.Init()

// Set a value without TTL
s.Set("foo", "bar", nil)

// Retrieve a value
val, ok := s.Get("foo")
if ok {
    fmt.Println(val) // "bar"
}

// Check if key exists
exists := s.Exists("foo") // true

// Delete a key
s.Delete("foo")
```

### With TTL (Time-To-Live)

```go
ttl := 5000 // 5 seconds in milliseconds
s.Set("temp", "data", &ttl)

// Value is automatically expired after 5 seconds
val, ok := s.Get("temp")
```

### Bulk Operations

```go
// Set multiple values
data := map[string]string{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
}
ttl := 10000 // 10 seconds
s.SetMany(data, &ttl)

// Get multiple values
values := s.GetMany([]string{"key1", "key2", "key3"})
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
