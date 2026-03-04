package store

import (
	"sync"
	"time"
)

// Package store provides a simple in-memory key/value store implementation
// with optional per-key TTL (time-to-live). TTLs are stored as absolute
// Unix milliseconds timestamps. A ttl value of -1 indicates no expiration.

// Store defines the minimal in-memory store operations used by callers.
// Implementations should be safe for concurrent use.
type Store interface {
	Init() error
	Set(string, string, *int) error
	Get(string) (string, bool)
	Delete(string) error
	Exists(string) bool
	SetMany(map[string]string, *int) error
	GetMany([]string) []string
	DeleteMany([]string) error
}

// Value holds a stored string and an expiration timestamp (Unix milliseconds).
// If ttl == -1 the value does not expire.
type Value struct {
	value string
	ttl   int64
}

// MemLocalStore is a simple in-memory implementation of Store. It uses a
// read-write mutex to guard access to the underlying map. Note: some methods
// acquire a write lock even for reads because they may remove expired keys
// lazily while servicing a read.
type MemLocalStore struct {
	mu          sync.RWMutex
	data        map[string]Value
	janitorStop chan struct{} // channel to signal janitor goroutine to exit
}

// Init initializes internal structures. Safe to call multiple times.
// It also starts the background janitor using a default interval of 1 second.
func (s *MemLocalStore) Init() error {
	s.mu.Lock()
	if s.data == nil {
		s.data = make(map[string]Value)
	}
	s.mu.Unlock()

	// start janitor if not already running
	s.StartJanitor(time.Second)
	return nil
}

// Set stores a single key with an optional TTL. The provided ttl is treated
// as an offset in milliseconds from now and converted to an absolute Unix
// millisecond timestamp stored in Value.ttl. A nil or zero ttl means no
// expiration (ttl == -1).
func (s *MemLocalStore) Set(k, v string, ttl *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]Value)
	}

	var ttlMili int64
	if ttl == nil || *ttl == 0 {
		// No expiration
		ttlMili = -1
	} else {
		ttlMili = time.Now().UnixMilli() + int64(*ttl)
	}
	s.data[k] = Value{value: v, ttl: ttlMili}
	return nil
}

// Get returns the value for a key and whether it existed and was not
// expired. Expired keys are removed lazily on access.
func (s *MemLocalStore) Get(k string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[k]
	if !ok {
		return "", false
	}

	// Lazy expiration: remove and report missing if expired.
	if val.ttl != -1 && time.Now().UnixMilli() > val.ttl {
		delete(s.data, k)
		return "", false
	}

	return val.value, true
}

// Delete removes a key from the store.
func (s *MemLocalStore) Delete(k string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, k)
	return nil
}

// Exists reports whether a key exists and has not expired. Expired keys are
// removed lazily.
func (s *MemLocalStore) Exists(k string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[k]
	if !ok {
		return false
	}

	// Lazy expiration: remove and report missing if expired.
	if val.ttl != -1 && time.Now().UnixMilli() > val.ttl {
		delete(s.data, k)
		return false
	}

	return true
}

// SetMany stores multiple keys with the same optional TTL. TTL semantics are
// the same as for Set: nil or zero = no expiration, otherwise ttl is treated
// as milliseconds offset from now.
func (s *MemLocalStore) SetMany(kv map[string]string, ttl *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]Value)
	}

	var ttlMili int64
	if ttl == nil || *ttl == 0 {
		// No expiration
		ttlMili = -1
	} else {
		ttlMili = time.Now().UnixMilli() + int64(*ttl)
	}

	for k, v := range kv {
		s.data[k] = Value{value: v, ttl: ttlMili}
	}
	return nil
}

// GetMany returns values for the provided keys in the same order. For keys
// that do not exist or have expired an empty string is returned in that slot.
// Expiration is still checked lazily in case janitor hasn't run yet.
func (s *MemLocalStore) GetMany(keys []string) []string {
	if len(keys) == 0 {
		return []string{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	result := []string{}
	now := time.Now().UnixMilli()

	for _, key := range keys {
		val, ok := s.data[key]
		if !ok {
			result = append(result, "")
			continue
		}

		// Lazy expiration: remove expired keys and return empty slot.
		if val.ttl != -1 && now > val.ttl {
			delete(s.data, key)
			result = append(result, "")
		} else {
			result = append(result, val.value)
		}
	}

	return result
}

// DeleteMany removes multiple keys from the store.
func (s *MemLocalStore) DeleteMany(keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		delete(s.data, key)
	}
	return nil
}

// cleanupExpired iterates over the map and deletes any keys whose TTL has
// passed. It acquires the write lock itself and is safe to call from the
// janitor goroutine or other internal callers.
func (s *MemLocalStore) cleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	for k, v := range s.data {
		if v.ttl != -1 && now > v.ttl {
			delete(s.data, k)
		}
	}
}

// StartJanitor begins a background goroutine that periodically cleans up
// expired keys. The interval parameter controls how often the cleanup runs;
// a non-positive value defaults to 1 second. It is safe to call multiple
// times; only one goroutine will be started. The janitor can be stopped with
// StopJanitor.
func (s *MemLocalStore) StartJanitor(interval time.Duration) {
	if interval <= 0 {
		interval = time.Second
	}

	s.mu.Lock()
	if s.janitorStop != nil {
		// already running
		s.mu.Unlock()
		return
	}
	stop := make(chan struct{})
	s.janitorStop = stop
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.cleanupExpired()
			case <-stop:
				return
			}
		}
	}()
}

// StopJanitor signals the background janitor goroutine to exit and waits for
// it to shut down by closing the stop channel. Calling StopJanitor on a store
// with no running janitor is a no-op.
func (s *MemLocalStore) StopJanitor() {
	s.mu.Lock()
	stop := s.janitorStop
	if stop == nil {
		s.mu.Unlock()
		return
	}
	// clear field while holding lock
	s.janitorStop = nil
	s.mu.Unlock()

	close(stop)
}
