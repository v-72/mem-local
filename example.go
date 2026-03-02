package main

import (
	"fmt"
	"time"

	"github.com/v-72/mem-local/store"
)

func main() {
	fmt.Println("=== mem-local Store Examples ===\n")

	// Initialize the store
	s := &store.MemLocalStore{}
	err := s.Init()
	if err != nil {
		fmt.Println("Error initializing store:", err)
		return
	}

	// Example 1: Basic Set and Get
	fmt.Println("1. Basic Set and Get")
	s.Set("name", "Alice", nil)
	val, ok := s.Get("name")
	fmt.Printf("   Get('name'): %v (found: %v)\n", val, ok)

	// Example 2: Set and Get with Exists
	fmt.Println("\n2. Using Exists to Check Keys")
	s.Set("city", "New York", nil)
	if s.Exists("city") {
		fmt.Println("   'city' key exists in store")
	}
	if !s.Exists("country") {
		fmt.Println("   'country' key does not exist in store")
	}

	// Example 3: Delete a Key
	fmt.Println("\n3. Delete a Key")
	s.Set("temp", "temporary", nil)
	fmt.Printf("   Before delete - Get('temp'): %v\n", val)
	s.Delete("temp")
	_, ok = s.Get("temp")
	fmt.Printf("   After delete - Get('temp') found: %v\n", ok)

	// Example 4: SetMany and GetMany (Bulk operations)
	fmt.Println("\n4. Bulk Operations with SetMany and GetMany")
	data := map[string]string{
		"user:1": "Alice",
		"user:2": "Bob",
		"user:3": "Charlie",
	}
	s.SetMany(data, nil)
	users := s.GetMany([]string{"user:1", "user:2", "user:3"})
	fmt.Println("   SetMany stored 3 users")
	fmt.Printf("   GetMany returns: %v\n", users)

	// Example 5: TTL - Keys with Expiration
	fmt.Println("\n5. Using TTL (Time-To-Live)")
	ttl := 2000 // 2 seconds
	s.Set("session:abc", "session_data", &ttl)
	val, ok = s.Get("session:abc")
	fmt.Printf("   Immediately after set - Get('session:abc'): %v (found: %v)\n", val, ok)

	fmt.Println("   Waiting 2.1 seconds for key to expire...")
	time.Sleep(2100 * time.Millisecond)

	val, ok = s.Get("session:abc")
	fmt.Printf("   After expiration - Get('session:abc'): %v (found: %v)\n", val, ok)

	// Example 6: TTL with SetMany
	fmt.Println("\n6. SetMany with TTL")
	cacheData := map[string]string{
		"cache:1": "data1",
		"cache:2": "data2",
	}
	cacheTTL := 1500 // 1.5 seconds
	s.SetMany(cacheData, &cacheTTL)
	cached := s.GetMany([]string{"cache:1", "cache:2"})
	fmt.Printf("   Cached data immediately: %v\n", cached)

	fmt.Println("   Waiting 1.6 seconds for cache to expire...")
	time.Sleep(1600 * time.Millisecond)

	expired := s.GetMany([]string{"cache:1", "cache:2"})
	fmt.Printf("   Cached data after expiration: %v\n", expired)

	// Example 7: Overwriting Keys
	fmt.Println("\n7. Overwriting Keys")
	s.Set("config", "version:1", nil)
	fmt.Printf("   Initial value: %v\n", mustGet(s, "config"))
	s.Set("config", "version:2", nil)
	fmt.Printf("   Updated value: %v\n", mustGet(s, "config"))

	// Example 8: Empty Values
	fmt.Println("\n8. Storing Empty Values")
	s.Set("empty", "", nil)
	empty, ok := s.Get("empty")
	fmt.Printf("   Set empty string - Get returns: '%v' (found: %v)\n", empty, ok)

	// Example 9: Mixed Existing and Non-existing Keys in GetMany
	fmt.Println("\n9. GetMany with Mixed Keys")
	s.Set("exists:1", "value1", nil)
	s.Set("exists:3", "value3", nil)
	mixed := s.GetMany([]string{"exists:1", "nonexistent", "exists:3"})
	fmt.Printf("   GetMany on mixed keys: %v\n", mixed)

	fmt.Println("\n=== End of Examples ===")
}

// Helper function to get value safely
func mustGet(s *store.MemLocalStore, key string) string {
	val, _ := s.Get(key)
	return val
}
