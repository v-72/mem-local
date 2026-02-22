package store

import (
	"testing"
)

func TestInit(t *testing.T) {
	s := &MemLocalStore{}
	err := s.Init()
	if err != nil {
		t.Errorf("Init() returned error: %v", err)
	}
	if s.Data == nil {
		t.Error("Init() did not initialize Data map")
	}
	if len(s.Data) != 0 {
		t.Errorf("Init() did not create empty map, got length %d", len(s.Data))
	}
}

func TestSet(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	result := s.Set("key1", "value1")
	if !result {
		t.Error("Set() returned false, expected true")
	}

	if s.Data["key1"] != "value1" {
		t.Errorf("Set() did not store value correctly, got %s", s.Data["key1"])
	}
}

func TestSetMultiple(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	s.Set("key1", "value1")
	s.Set("key2", "value2")
	s.Set("key3", "value3")

	if len(s.Data) != 3 {
		t.Errorf("Set() multiple values failed, expected 3 items, got %d", len(s.Data))
	}
}

func TestSetOverwrite(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	s.Set("key1", "value1")
	s.Set("key1", "value2")

	if s.Data["key1"] != "value2" {
		t.Errorf("Set() overwrite failed, expected value2, got %s", s.Data["key1"])
	}
}

func TestGet(t *testing.T) {
	s := &MemLocalStore{}
	s.Init()

	s.Set("key1", "value1")
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

	result := s.Set("key1", "value1")

	if !result {
		t.Error("Set() returned false, expected true")
	}
	if s.Data == nil {
		t.Error("Set() did not initialize Data map when nil")
	}
	if s.Data["key1"] != "value1" {
		t.Errorf("Set() without Init failed, got %s", s.Data["key1"])
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
		s.Set(k, v)
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

	s.Set("emptyKey", "")
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
