# mem-local

A small Go project providing a simple in-memory local store. This repository contains a minimal example application and a lightweight `store` implementation intended for local, ephemeral usage and testing.

**Status:** minimal example / reference implementation

**Prerequisites:**
- Go 1.20+

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
- [main.go](main.go): example program / entrypoint
- [store/mem_local.go](store/mem_local.go): in-memory store implementation
- [go.mod](go.mod): Go module definition

**Usage Notes**
- The store is intentionally in-memory: data is not persisted across restarts.
- Designed for testing, prototypes, or local caches where persistence is unnecessary.

**Contributing**
- Feel free to open issues or PRs. Keep changes small and focused.

**License**
- No license included. Add one (`LICENSE`) if you intend to open-source this project.
# mem-local
local cache
