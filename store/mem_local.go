package store

import (
	"sync"
	"time"
)

type Store interface {
	Init() error
	Set(string, string) error
	Get(string) (string, bool)
	Delete(string) error
	Exists(string) bool
	SetMany(map[string]string) error
	GetMany([]string) []string
	DeleteMany([]string) error
}

type Value struct {
	value string
	ttl   int64
}

type MemLocalStore struct {
	mu   sync.RWMutex
	data map[string]Value
}

func (s *MemLocalStore) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]Value)
	return nil
}

func (s *MemLocalStore) Set(k, v string, ttl *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]Value)
	}

	ttlMili := time.Now().UnixMilli() + int64(getTTL(ttl))
	s.data[k] = Value{value: v, ttl: ttlMili}
	return nil
}

func (s *MemLocalStore) Get(k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[k]
	return val.value, ok
}

func (s *MemLocalStore) Delete(k string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, k)
	return nil
}

func (s *MemLocalStore) Exists(k string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[k]
	return ok
}

func (s *MemLocalStore) SetMany(kv map[string]string, ttl *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]Value)
	}
	ttlMili := time.Now().UnixMilli() + int64(getTTL(ttl))
	for k, v := range kv {
		s.data[k] = Value{value: v, ttl: ttlMili}
	}
	return nil
}

func (s *MemLocalStore) GetMany(keys []string) []string {
	if len(keys) == 0 {
		return []string{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	result := []string{}

	for _, key := range keys {
		result = append(result, s.data[key].value)
	}

	return result
}

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

func getTTL(ttl *int) int {
	if ttl == nil || *ttl == 0 {
		return 0
	}
	return *ttl
}
