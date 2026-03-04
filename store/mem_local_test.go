package store

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestConcurrentSetMultipleKeys(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			val := "value" + strconv.Itoa(i)
			if err := s.Set(key, val, nil); err != nil {
				t.Errorf("Set failed for %s: %v", key, err)
			}
		}()
	}
	wg.Wait()

	// Check all keys
	for i := 0; i < n; i++ {
		key := "key" + strconv.Itoa(i)
		val, ok := s.Get(key)
		if !ok || val != "value"+strconv.Itoa(i) {
			t.Errorf("Get(%s) = %v, %v; want value%d, true", key, val, ok, i)
		}
	}
}

func TestConcurrentSetAndGetMultipleKeys(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	const n = 100
	var wg sync.WaitGroup

	// Set keys concurrently
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			val := "value" + strconv.Itoa(i)
			if err := s.Set(key, val, nil); err != nil {
				t.Errorf("Set failed for %s: %v", key, err)
			}
		}()
	}

	// Get keys concurrently
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			val, ok := s.Get(key)
			// Accept either not found or correct value, since gets may race with sets
			if ok && val != "value"+strconv.Itoa(i) {
				t.Errorf("Concurrent Get(%s) = %v, want value%d or not found", key, val, i)
			}
		}()
	}

	wg.Wait()
}

func TestInit(t *testing.T) {
	s := &MemLocalStore{}
	err := s.Init()
	if err != nil {
		t.Errorf("Init() returned error: %v", err)
	}
	// After Init, getting any key should return not-ok.
	if val, ok := s.Get("any"); ok || val != "" {
		t.Errorf("Init() did not create empty store, got %v, %v", val, ok)
	}
}

func TestSet(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Errorf("Set() returned error: %v", err)
	}

	val, ok := s.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Set() did not store value correctly, got %v, %v", val, ok)
	}
}

func TestSetMultiple(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	keys := []struct{ k, v string }{{"key1", "value1"}, {"key2", "value2"}, {"key3", "value3"}}
	for _, kv := range keys {
		if err := s.Set(kv.k, kv.v, nil); err != nil {
			t.Fatalf("Set(%s) error: %v", kv.k, err)
		}
	}

	for _, kv := range keys {
		val, ok := s.Get(kv.k)
		if !ok || val != kv.v {
			t.Errorf("Get(%s) returned %v,%v expected %s,true", kv.k, val, ok, kv.v)
		}
	}
}

func TestSetOverwrite(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := s.Set("key1", "value2", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	if val, ok := s.Get("key1"); !ok || val != "value2" {
		t.Errorf("Set() overwrite failed, expected value2, got %v,%v", val, ok)
	}
}

func TestGet(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	val, ok := s.Get("key1")

	if !ok {
		t.Error("Get() returned false for existing key, expected true")
	}
	if val != "value1" {
		t.Errorf("Get() returned wrong value, expected value1, got %s", val)
	}
}

func TestGetNonExistent(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	val, ok := s.Get("nonexistent")

	if ok {
		t.Error("Get() returned true for non-existent key, expected false")
	}
	if val != "" {
		t.Errorf("Get() returned non-empty value for non-existent key, got %s", val)
	}
}

func TestGetWithoutInit(t *testing.T) {
	s := &MemLocalStore{}
	val, ok := s.Get("key1")

	if ok {
		t.Error("Get() returned true for uninitialized store, expected false")
	}
	if val != "" {
		t.Errorf("Get() returned non-empty value, got %s", val)
	}
}

func TestSetWithoutInit(t *testing.T) {
	s := &MemLocalStore{}
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Errorf("Set() returned error, expected nil: %v", err)
	}
	val, ok := s.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Set() without Init failed, got %v,%v", val, ok)
	}
}

func TestGetMultiple(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	testCases := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range testCases {
		if err := s.Set(k, v, nil); err != nil {
			t.Fatalf("Set(%s) error: %v", k, err)
		}
	}

	for k, expected := range testCases {
		val, ok := s.Get(k)
		if !ok {
			t.Errorf("Get(%s) returned false, expected true", k)
		}
		if val != expected {
			t.Errorf("Get(%s) returned %s, expected %s", k, val, expected)
		}
	}
}

func TestEmptyStringValue(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("emptyKey", "", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	val, ok := s.Get("emptyKey")

	if !ok {
		t.Error("Get() returned false for key with empty string value, expected true")
	}
	if val != "" {
		t.Errorf("Get() returned non-empty value for empty string key, got %s", val)
	}
}

// func TestStoreInterface(t *testing.T) {
// 	var _ Store = (*MemLocalStore)(nil)
// }

func TestDelete(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	val, ok := s.Get("key1")

	if !ok {
		t.Error("Get() returned false for existing key, expected true")
	}
	if val != "value1" {
		t.Errorf("Get() returned wrong value, expected value1, got %s", val)
	}

	if err := s.Delete("key1"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	val, ok = s.Get("key1")
	if ok {
		t.Error("Get() returned true for deleted key, expected false")
	}
	if val != "" {
		t.Errorf("Get() returned non-empty value for deleted key, got %s", val)
	}
}

func TestExists(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if !s.Exists("key1") {
		t.Error("Exists() returned false for existing key, expected true")
	}
	if s.Exists("nonexistent") {
		t.Error("Exists() returned true for non-existent key, expected false")
	}
}

func TestSetMany(t *testing.T) {
	s := &MemLocalStore{}
	testCases := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	if err := s.SetMany(testCases, nil); err != nil {
		t.Fatalf("SetMany() error: %v", err)
	}

	for k, expected := range testCases {
		val, ok := s.Get(k)
		if !ok {
			t.Errorf("Get(%s) returned false, expected true", k)
		}
		if val != expected {
			t.Errorf("Get(%s) returned %s, expected %s", k, val, expected)
		}
	}
}

func TestGetMany(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	testCases := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	s.SetMany(testCases, nil)

	keys := []string{"key1", "key2", "key3"}
	values := s.GetMany(keys)

	if len(values) != len(keys) {
		t.Errorf("GetMany() returned %d values, expected %d", len(values), len(keys))
	}

	for i, key := range keys {
		if values[i] != testCases[key] {
			t.Errorf("GetMany() returned %s for key %s, expected %s", values[i], key, testCases[key])
		}
	}
}

func TestGetManyEmpty(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	keys := []string{"key1", "key2", "key3"}
	values := s.GetMany(keys)

	if len(values) != len(keys) {
		t.Errorf("GetMany() returned %d values, expected %d", len(values), len(keys))
	}

	for _, value := range values {
		if value != "" {
			t.Errorf("GetMany() returned non-empty value %s, expected empty", value)
		}
	}
}

func TestGetManyNonExistent(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	keys := []string{"nonexistent1", "nonexistent2", "nonexistent3"}
	values := s.GetMany(keys)

	if len(values) != len(keys) {
		t.Errorf("GetMany() returned %d values, expected %d", len(values), len(keys))
	}

	for _, value := range values {
		if value != "" {
			t.Errorf("GetMany() returned non-empty value %s, expected empty", value)
		}
	}
}

func TestGetManyMixed(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	keys := []string{"key1", "nonexistent", "key2"}
	values := s.GetMany(keys)

	if len(values) != len(keys) {
		t.Errorf("GetMany() returned %d values, expected %d", len(values), len(keys))
	}

	if values[0] != "value1" {
		t.Errorf("GetMany() returned %s for key %s, expected %s", values[0], keys[0], "value1")
	}

	if values[1] != "" {
		t.Errorf("GetMany() returned %s for key %s, expected %s", values[1], keys[1], "")
	}

	if values[2] != "" {
		t.Errorf("GetMany() returned %s for key %s, expected %s", values[2], keys[2], "")
	}
}

func TestGetManyEmptyKeys(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	keys := []string{}
	values := s.GetMany(keys)

	if len(values) != 0 {
		t.Errorf("GetMany() returned %d values, expected 0", len(values))
	}
}

func TestGetManyWithoutInit(t *testing.T) {
	s := &MemLocalStore{}

	keys := []string{"key1", "key2"}
	values := s.GetMany(keys)

	if len(values) != len(keys) {
		t.Errorf("GetMany() returned %d values, expected %d", len(values), len(keys))
	}

	for _, value := range values {
		if value != "" {
			t.Errorf("GetMany() returned non-empty value %s, expected empty", value)
		}
	}
}

func TestGetManyConcurrent(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	const n = 100
	var wg sync.WaitGroup

	kv := make(map[string]string)
	for i := 0; i < n; i++ {
		i := i
		kv["key"+strconv.Itoa(i)] = "value" + strconv.Itoa(i)
	}
	s.SetMany(kv, nil)

	// Get keys concurrently
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			values := s.GetMany([]string{key})
			if len(values) != 1 || values[0] != "value"+strconv.Itoa(i) {
				t.Errorf("GetMany() returned %v for key %s, expected [value%d]", values, key, i)
			}
		}()
	}

	wg.Wait()
}

func TestDeleteMany(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Add some test data
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := s.Set("key2", "value2", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := s.Set("key3", "value3", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Test deleting multiple keys
	keysToDelete := []string{"key1", "key3"}
	if err := s.DeleteMany(keysToDelete); err != nil {
		t.Fatalf("DeleteMany error: %v", err)
	}

	// Verify deleted keys are gone
	for _, key := range keysToDelete {
		_, ok := s.Get(key)
		if ok {
			t.Errorf("Key %s should have been deleted but still exists", key)
		}
	}

	// Verify remaining key still exists
	val, ok := s.Get("key2")
	if !ok {
		t.Errorf("Key key2 should still exist")
	}
	if val != "value2" {
		t.Errorf("Key key2 has wrong value: %s, expected value2", val)
	}
}

func TestDeleteManyNonExistentKeys(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Add some test data
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Test deleting non-existent keys
	keysToDelete := []string{"nonexistent1", "nonexistent2"}
	if err := s.DeleteMany(keysToDelete); err != nil {
		t.Fatalf("DeleteMany error: %v", err)
	}

	// Verify existing key still exists
	val, ok := s.Get("key1")
	if !ok {
		t.Errorf("Key key1 should still exist")
	}
	if val != "value1" {
		t.Errorf("Key key1 has wrong value: %s, expected value1", val)
	}
}

func TestDeleteManyEmptyKeys(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Test deleting with empty keys slice
	keysToDelete := []string{}
	if err := s.DeleteMany(keysToDelete); err != nil {
		t.Fatalf("DeleteMany error: %v", err)
	}

	// Verify store is unchanged
	val, ok := s.Get("key1")
	if ok {
		t.Errorf("Key key1 should not exist: %s", val)
	}
}

func TestDeleteManyWithoutInit(t *testing.T) {
	s := &MemLocalStore{}

	// Test deleting without initializing the store
	keysToDelete := []string{"key1", "key2"}
	if err := s.DeleteMany(keysToDelete); err != nil {
		t.Fatalf("DeleteMany error: %v", err)
	}

	// Verify store is still empty
	val, ok := s.Get("key1")
	if ok {
		t.Errorf("Key key1 should not exist: %s", val)
	}
}

func TestDeleteManyConcurrent(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	const n = 100
	var wg sync.WaitGroup

	// Add test data
	for i := 0; i < n; i++ {
		key := "key" + strconv.Itoa(i)
		if err := s.Set(key, "value"+strconv.Itoa(i), nil); err != nil {
			t.Fatalf("Set error: %v", err)
		}
	}

	// Delete keys concurrently
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			if err := s.DeleteMany([]string{key}); err != nil {
				t.Errorf("DeleteMany error for key %s: %v", key, err)
			}
		}()
	}

	wg.Wait()

	// Verify all keys are deleted
	for i := 0; i < n; i++ {
		key := "key" + strconv.Itoa(i)
		_, ok := s.Get(key)
		if ok {
			t.Errorf("Key %s should have been deleted but still exists", key)
		}
	}
}
func TestSetWithTTL(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	ttl := 5000 // 5 seconds
	if err := s.Set("key1", "value1", &ttl); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	val, ok := s.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Set with TTL failed, got %v, %v", val, ok)
	}
}

func TestGetWithTTLExpired(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Set with very short TTL (1 millisecond)
	ttl := 1
	if err := s.Set("key1", "value1", &ttl); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Wait a bit to ensure expiration
	time.Sleep(10 * time.Millisecond)

	val, ok := s.Get("key1")
	if ok {
		t.Errorf("Get should return false for expired key, got %v, %v", val, ok)
	}
	if val != "" {
		t.Errorf("Get should return empty string for expired key, got %s", val)
	}
}

func TestGetWithTTLNotExpired(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Set with long TTL
	ttl := 10000 // 10 seconds
	if err := s.Set("key1", "value1", &ttl); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	val, ok := s.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Get should return value for non-expired key, got %v, %v", val, ok)
	}
}

func TestExistsWithTTLExpired(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Set with very short TTL
	ttl := 1
	if err := s.Set("key1", "value1", &ttl); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	if !s.Exists("key1") {
		t.Error("Exists should return true before expiration")
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	if s.Exists("key1") {
		t.Error("Exists should return false after expiration")
	}
}

func TestGetManyWithTTLExpired(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Set some keys without TTL
	if err := s.Set("key1", "value1", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Set key with very short TTL
	ttl := 1
	if err := s.Set("key2", "value2", &ttl); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Set another key without TTL
	if err := s.Set("key3", "value3", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	values := s.GetMany([]string{"key1", "key2", "key3"})

	if len(values) != 3 {
		t.Errorf("GetMany returned %d values, expected 3", len(values))
	}

	if values[0] != "value1" {
		t.Errorf("GetMany: key1 should be 'value1', got '%s'", values[0])
	}

	if values[1] != "" {
		t.Errorf("GetMany: key2 should be empty (expired), got '%s'", values[1])
	}

	if values[2] != "value3" {
		t.Errorf("GetMany: key3 should be 'value3', got '%s'", values[2])
	}
}

func TestLazyExpirationRemovesKey(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// Set with very short TTL
	ttl := 1
	if err := s.Set("key1", "value1", &ttl); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Get should remove the expired key
	_, _ = s.Get("key1")

	// Now insert new key and verify lazy deletion removed the old one
	if err := s.Set("newkey", "newvalue", nil); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Try to get the original expired key - it should not exist
	val, ok := s.Get("key1")
	if ok {
		t.Errorf("Lazy expiration should have removed key1, but got %v, %v", val, ok)
	}
}

func TestSetManyWithTTL(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	ttl := 5000 // 5 seconds
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	if err := s.SetMany(data, &ttl); err != nil {
		t.Fatalf("SetMany error: %v", err)
	}

	for k, expected := range data {
		val, ok := s.Get(k)
		if !ok || val != expected {
			t.Errorf("SetMany with TTL: Get(%s) = %v, %v; expected %s, true", k, val, ok, expected)
		}
	}
}

func TestTTLZeroMeansNoExpiration(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	// TTL of 0 should mean no expiration
	ttl := 0
	if err := s.Set("key1", "value1", &ttl); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Even after a long wait, key should still exist because ttl is 0
	time.Sleep(10 * time.Millisecond)

	val, ok := s.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("TTL of 0 should mean no expiration, got %v, %v", val, ok)
	}
}

func TestBackgroundJanitorRemovesExpired(t *testing.T) {
s := &MemLocalStore{}
s.Init()

// restart janitor with a fast interval so the test completes quickly
s.StopJanitor()
s.StartJanitor(10 * time.Millisecond)
defer s.StopJanitor()

ttl := 1 // 1 millisecond
if err := s.Set("key1", "value1", &ttl); err != nil {
t.Fatalf("Set error: %v", err)
}

// sleep long enough for key to expire and janitor to run
time.Sleep(50 * time.Millisecond)

if s.Exists("key1") {
t.Error("Background janitor should have removed expired key1")
}
}

func TestJanitorStoppedDoesNotPanic(t *testing.T) {
s := &MemLocalStore{}
s.Init()
// ensure StopJanitor is safe when called multiple times
s.StopJanitor()
s.StopJanitor()
}
