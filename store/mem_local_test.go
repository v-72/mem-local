package store

import (
	"strconv"
	"sync"
	"testing"
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
			if err := s.Set(key, val); err != nil {
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
			if err := s.Set(key, val); err != nil {
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
	if err := s.Set("key1", "value1"); err != nil {
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
		if err := s.Set(kv.k, kv.v); err != nil {
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
	if err := s.Set("key1", "value1"); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := s.Set("key1", "value2"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	if val, ok := s.Get("key1"); !ok || val != "value2" {
		t.Errorf("Set() overwrite failed, expected value2, got %v,%v", val, ok)
	}
}

func TestGet(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("key1", "value1"); err != nil {
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
	if err := s.Set("key1", "value1"); err != nil {
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
		if err := s.Set(k, v); err != nil {
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
	if err := s.Set("emptyKey", ""); err != nil {
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

func TestStoreInterface(t *testing.T) {
	var _ Store = (*MemLocalStore)(nil)
}

func TestDelete(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()
	if err := s.Set("key1", "value1"); err != nil {
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
