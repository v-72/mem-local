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
- `Set(key, value string, ttl *int) error` — Store a value with optional TTL in milliseconds. Pass `nil` or `0` for no expiration.
- `Get(key string) (string, bool)` — Retrieve a value. Checks for expiration and removes expired keys (lazy expiration).
- `Delete(key string) error` — Remove a key
- `Exists(key string) bool` — Check if a key exists. Checks for expiration and removes expired keys (lazy expiration).
- `SetMany(kv map[string]string, ttl *int) error` — Store multiple values with optional TTL
- `GetMany(keys []string) []string` — Retrieve multiple values. Checks expiration for each key (lazy expiration).
- `StartJanitor(interval time.Duration)` — Start a background goroutine that periodically purges expired entries. Called automatically by `Init()` with a 1‑second interval; you can call again with a custom interval (e.g. for faster cleanup in tests).
- `StopJanitor()` — Stops the background janitor, if running (safe to call multiple times).

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
// Set a value with TTL of 5 seconds
ttl := 5000 // in milliseconds
s.Set("temp", "data", &ttl)

// Value is accessible before timeout
val, ok := s.Get("temp") // returns "data", true

// After 5 seconds, accessing the key will return empty and remove it
// (lazy expiration - checked on access)
time.Sleep(5100 * time.Millisecond)
val, ok = s.Get("temp") // returns "", false
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
- **Lazy Expiration:** Keys with TTL are checked for expiration when accessed via `Get()`, `Exists()`, or `GetMany()` (lazy removal). Additionally, a background janitor goroutine—started automatically by `Init()`—periodically scans the store and deletes expired entries, ensuring they disappear even if they are never touched again.
- Pass `nil` or `0` as TTL to store keys indefinitely with no expiration.
- All operations are thread-safe and can be used concurrently.

## Testing
Run the test suite with:

```bash
go test ./store -v
```

The project includes 35+ test cases covering:
- Concurrent operations
- TTL and lazy expiration behavior
- Edge cases (empty values, non-existent keys, uninitialized store)
- Bulk operations

## Contributing
- Issues and PRs are welcome. Please keep changes small and focused.

## License
- No license included. Add one (`LICENSE`) if you intend to open-source this project.
